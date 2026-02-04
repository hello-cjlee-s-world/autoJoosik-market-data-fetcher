package scheduler

import (
	"autoJoosik-market-data-fetcher/internal/autoSellerService"
	"autoJoosik-market-data-fetcher/internal/kiwoomApi"
	"autoJoosik-market-data-fetcher/internal/model"
	"autoJoosik-market-data-fetcher/internal/repository"
	"autoJoosik-market-data-fetcher/internal/utils"
	"autoJoosik-market-data-fetcher/pkg/logger"
	"context"
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/jackc/pgx/v5/pgxpool"
)

// 작업 시그니처와 레지스트리 (task_type -> 함수)
type TaskFunc func(ctx context.Context) error
type Registry map[string]TaskFunc

// DB에서 읽어올 스케줄 설정
type TaskConfig struct {
	Name     string
	Schedule string // cron 또는 "every 10s", "@every 5m"
	TaskType string
}

// 스케줄러 러너
type Runner struct {
	pool     *pgxpool.Pool
	registry Registry
	s        *gocron.Scheduler
	running  atomic.Bool
	enabled  atomic.Bool
}

// 거래 목록
var TradeUniverses []model.TbTradeUniverseEntity
var StkCdList []string

// Runner 생성
func NewRunner(pool *pgxpool.Pool, reg Registry, runner *Runner) *Runner {
	s := gocron.NewScheduler(time.Local)
	s.TagsUnique()
	runner.pool = pool
	runner.registry = reg
	runner.s = s
	runner.enabled.Store(true)
	return runner
}

var ErrRunnerNotInitialized = errors.New("runner not initialized")

func (r *Runner) Initialized() bool {
	return r != nil && r.s != nil && r.pool != nil && r.registry != nil
}

func (r *Runner) Start(ctx context.Context) error {
	if !r.Initialized() {
		return ErrRunnerNotInitialized
	}
	if !r.enabled.Load() {
		return nil
	}
	if !r.running.CompareAndSwap(false, true) {
		return nil
	}
	if err := r.Reload(ctx); err != nil {
		r.running.Store(false)
		return err
	}
	r.s.StartAsync()
	return nil
}

func (r *Runner) Stop() error {
	if !r.Initialized() {
		return ErrRunnerNotInitialized
	}
	if !r.running.CompareAndSwap(true, false) {
		return nil
	}
	r.s.Stop()
	return nil
}
func (r *Runner) Enable(ctx context.Context) error {
	if !r.Initialized() {
		return ErrRunnerNotInitialized
	}
	r.enabled.Store(true)
	return r.Start(ctx)
}

func (r *Runner) Disable() error {
	if !r.Initialized() {
		return ErrRunnerNotInitialized
	}
	r.enabled.Store(false)
	return r.Stop()
}

func (r *Runner) Reload(ctx context.Context) error {
	if !r.Initialized() {
		return ErrRunnerNotInitialized
	}
	r.s.Clear()

	cfgs, err := LoadTaskConfigs(ctx, r.pool)
	if err != nil {
		return err
	}
	for _, c := range cfgs {
		if err := r.scheduleOne(ctx, c); err != nil {
			return err
		}
	}
	return nil
}

// DB에서 스케줄 읽기
func LoadTaskConfigs(ctx context.Context, pool *pgxpool.Pool) ([]TaskConfig, error) {
	// enabled 컬럼 없을 수 있으니 심플하게 3개만
	rows, err := pool.Query(ctx, `SELECT name, schedule, task_type FROM tb_schedule_info WHERE enabled = true`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []TaskConfig
	for rows.Next() {
		var c TaskConfig
		if err := rows.Scan(&c.Name, &c.Schedule, &c.TaskType); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// 하나의 잡을 스케줄링
func (r *Runner) scheduleOne(ctx context.Context, c TaskConfig) error {
	task, ok := r.registry[c.TaskType]
	if !ok {
		return errors.New("unknown task_type: " + c.TaskType)
	}

	spec := strings.TrimSpace(c.Schedule)
	specLower := strings.ToLower(spec)

	// 1) "every 10s" or "@every 10s"
	if strings.HasPrefix(specLower, "every ") || strings.HasPrefix(specLower, "@every ") {
		d := strings.TrimSpace(strings.TrimPrefix(specLower, "@every "))
		d = strings.TrimSpace(strings.TrimPrefix(d, "every "))
		n, unit, err := parseEverySpec(d) // "10s" -> (10, "s")
		if err != nil {
			return err
		}
		var j *gocron.Job
		switch unit {
		case "s":
			j, err = r.s.Every(n).Seconds().Do(func() { _ = task(ctx) })
		case "m":
			j, err = r.s.Every(n).Minutes().Do(func() { _ = task(ctx) })
		case "h":
			j, err = r.s.Every(n).Hours().Do(func() { _ = task(ctx) })
		default:
			return errors.New("unsupported every unit: " + unit)
		}
		if err != nil {
			return err
		}
		j.Tag(c.Name, c.TaskType)
		return nil
	}

	// 2) CRON 표현 (기본/초 포함 둘 다 시도)
	j, err := r.s.Cron(spec).Do(func() { _ = task(ctx) })
	if err != nil {
		j, err = r.s.CronWithSeconds(spec).Do(func() { _ = task(ctx) })
		if err != nil {
			return err
		}
	}
	j.Tag(c.Name, c.TaskType)
	return nil
}

// "10s", "5m", "2h" 파싱
func parseEverySpec(s string) (int, string, error) {
	if len(s) < 2 {
		return 0, "", errors.New("invalid every spec: " + s)
	}
	unit := s[len(s)-1:] // 마지막 글자
	numStr := strings.TrimSpace(s[:len(s)-1])
	n, err := strconv.Atoi(numStr)
	if err != nil || n <= 0 {
		return 0, "", errors.New("invalid every number: " + numStr)
	}
	switch unit {
	case "s", "m", "h":
		return n, unit, nil
	default:
		return 0, "", errors.New("invalid every unit (use s/m/h): " + unit)
	}
}

func GetSchedule(ctx context.Context, pool *pgxpool.Pool, r *Runner) {
	// 실제 업무 함수들 등록 (task_type -> 함수). 필요 시 pool 캡쳐해서 사용
	reg := Registry{
		"GetTradeUniverse": func(ctx context.Context) error {
			rst, err := repository.GetTradeUniverse(ctx, pool)
			if err != nil {
				return err
			}
			if len(rst) > 0 {
				TradeUniverses = rst
				// 누적 append 방지
				StkCdList = StkCdList[:0]
				for _, u := range TradeUniverses {
					StkCdList = append(StkCdList, u.StkCd)
				}

			}
			return nil
		},
		"GetTradeInfoLog": func(ctx context.Context) error {
			// fix !! 가지고 있는 종목 별로 바꿔야됨
			var entList []model.TbTradeInfoLogEntity
			for _, stkCd := range StkCdList {
				rst, err := kiwoomApi.GetTradeInfoLog(stkCd)
				if err != nil {
					return err
				}
				ToEntList := model.ToTbTradeInfoLogEntity(rst, stkCd)
				entList = append(entList, ToEntList...)
			}

			err := repository.UpsertTradeInfoBatch(ctx, pool, entList)
			if err != nil {
				return err
			}
			return nil
		},
		"UpsertStockInfo": func(ctx context.Context) error {
			var entList []model.TbStockInfoEntity
			for _, stkCd := range StkCdList {
				rst, err := kiwoomApi.GetStockInfo(stkCd)
				if err != nil {
					return err
				}
				ent := model.ToTbStockInfoEntity(rst)
				entList = append(entList, ent)
			}

			err := repository.UpsertStockInfoBatch(ctx, pool, entList)
			if err != nil {
				return err
			}
			return nil
		},
		"SellOrBuy": func(ctx context.Context) error {
			return autoSellerService.DecideAndExecute(ctx, pool)
		},
		"CalStockScore": func(ctx context.Context) error {
			var entList []model.TbStockScoreEntity
			for _, stkCd := range StkCdList {
				bullBearEntity, err := repository.GetBullBearValue(ctx, pool, stkCd)
				if err != nil {
					return err
				}
				tbStockInfoEntity, err := repository.GetStockFundamental(ctx, pool, stkCd)
				if err != nil {
					return err
				}
				ent, err := calcScoreToEntity(bullBearEntity, tbStockInfoEntity, stkCd)
				if err != nil {
					return err
				}
				entList = append(entList, ent)
			}

			return repository.UpsertStockScoreBatch(ctx, pool, entList)
		},
		"CalAssetStats": func(ctx context.Context) error {
			accountId := int64(0)
			account, _ := repository.GetVirtualAccount(ctx, pool, accountId)
			assetList, _ := repository.GetHoldingPositions(ctx, pool, accountId)

			var cash = int(account.CashBalance)
			var stockValue = 0
			for _, asset := range assetList {
				stockValue += int(asset.Qty) * int(asset.AvgPrice)
			}
			var total = cash + stockValue

			err := repository.InsertVirtualAssetDaily(ctx, pool, total, stockValue, cash)
			if err != nil {
				return err
			}

			return nil
		},
	}

	// 러너 생성 & 시작
	NewRunner(pool, reg, r)

	// 장 상태 감시: 30초마다 체크
	guard := time.NewTicker(30 * time.Second)
	defer guard.Stop()

	// schedule_info 리로드용(장 열렸을 때만 수행 추천)
	reloadTicker := time.NewTicker(30 * time.Second)
	defer reloadTicker.Stop()

	if err := r.Start(ctx); err != nil {
		logger.Warn("scheduler start:", err)
	}
	logger.Info("[scheduler] started")

	for {
		select {
		case <-ctx.Done():
			r.Stop()
			logger.Info("[scheduler] stopped")
			return

		case <-guard.C:
			if !r.IsEnabled() {
				was := r.IsRunning()
				r.Stop()
				if was && !r.IsRunning() {
					logger.Info("[scheduler] stopped (manually disabled)")
				}
				break
			}
			now := time.Now()
			tradable := IsTradableTime(now)

			if tradable {
				was := r.IsRunning()
				if err := r.Start(ctx); err != nil {
					logger.Info("[scheduler] start error:", err)
					break
				}
				if !was && r.IsRunning() {
					logger.Info("[scheduler] started (market open)")
				}
			} else {
				was := r.IsRunning()
				r.Stop()
				if was && !r.IsRunning() {
					logger.Info("[scheduler] stopped (market closed)")
				}
			}

		case <-reloadTicker.C:
			if !r.IsEnabled() {
				break
			}
			if err := r.Reload(ctx); err != nil {
				logger.Info("[scheduler] reload error:", err)
			}
		}
	}
}

func calcScoreToEntity(
	bullEntity model.BullBearEntity, // R1,R2,R3,Volatility,LastPrice 같은 값 있다고 가정
	infoEntity model.TbStockInfoEntity, // Per,Pbr,Roe,Eps,ForExhRt,Cap 등이 string으로 들어있다고 가정
	stkCd string,
) (model.TbStockScoreEntity, error) {
	var ent model.TbStockScoreEntity
	// ===== 모멘텀 점수 =====
	momentum := 0.0
	if bullEntity.R1 > 0 {
		momentum += 10
	}
	if bullEntity.R2 > 0 {
		momentum += 15
	}
	if bullEntity.R3 > 0 {
		momentum += 25
	}

	// ===== 리스크 감점 =====
	risk := bullEntity.Volatility * 10
	risk = math.Min(risk, 30)

	// ===== 재무 점수 =====
	per, err := utils.ParseSignedFloat(infoEntity.Per)
	if err != nil {
		return model.TbStockScoreEntity{}, err
	}

	pbr, err := utils.ParseSignedFloat(infoEntity.Pbr)
	if err != nil {
		return model.TbStockScoreEntity{}, err
	}

	roe, err := utils.ParseSignedFloat(infoEntity.Roe)
	if err != nil {
		return model.TbStockScoreEntity{}, err
	}

	eps, err := utils.ParseSignedFloat(infoEntity.Eps)
	if err != nil {
		return model.TbStockScoreEntity{}, err
	}

	forExhRt, err := utils.ParseSignedFloat(infoEntity.ForExhRt)
	if err != nil {
		return model.TbStockScoreEntity{}, err
	}

	capVal, err := utils.ParseSignedFloat(infoEntity.Cap)
	if err != nil {
		return model.TbStockScoreEntity{}, err
	}

	lastPrice, err := utils.ParseSignedFloat(infoEntity.CurPrc)
	if err != nil {
		return model.TbStockScoreEntity{}, err
	}

	fund := 0.0
	if per > 0 && per < 15 {
		fund += 15
	}
	if pbr > 0 && pbr < 1.5 {
		fund += 10
	}
	if roe >= 10 {
		fund += 10
	}
	if eps > 0 {
		fund += 5
	}
	if forExhRt >= 10 {
		fund += 5
	}
	if capVal > 0 && capVal < 1_000_000_000_000 {
		fund -= 5
	} // 너무 소형주

	// ===== 최종 점수 =====
	scoreTotal := fund + momentum - risk
	scoreTotal = math.Max(0, math.Min(100, scoreTotal))

	now := time.Now()

	metaMap := map[string]any{
		"per":      infoEntity.Per,
		"pbr":      infoEntity.Pbr,
		"roe":      infoEntity.Roe,
		"eps":      infoEntity.Eps,
		"forExhRt": infoEntity.ForExhRt,
		"cap":      infoEntity.Cap,
	}
	metaBytes, err := json.Marshal(metaMap)
	if err != nil {
		logger.Error("CalcScore :: meta marshal error :: " + err.Error())
		return ent, err
	}

	// ===== 엔티티 매핑 =====
	ent = model.TbStockScoreEntity{
		StkCd:            stkCd, // 또는 bullEntity.StkCd
		ScoreTotal:       scoreTotal,
		ScoreFundamental: fund,
		ScoreMomentum:    momentum,
		ScoreMarket:      0,         // 아직 미사용이면 0
		ScoreRisk:        risk,      // "감점값" 그대로 저장 (원하면 -risk로 저장해도 됨)
		LastPrice:        lastPrice, // 없으면 infoEntity.CurPrc 파싱해서 넣어도 됨
		R1:               bullEntity.R1,
		R2:               bullEntity.R2,
		R3:               bullEntity.R3,
		Volatility:       bullEntity.Volatility,
		AsofTm:           now,
		Meta:             string(metaBytes),
		CreatedAt:        now, // UPSERT면 사실상 Update에서만 의미. Insert 시에만 넣고 싶으면 repo에서 처리해도 됨
		UpdatedAt:        now,
	}
	logger.Debug("calcScoreToEntity :: Success :: " + ent.StkCd)
	return ent, nil
}

// 거래 시간,날짜 제어
func IsTradableTime(now time.Time) bool {
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		return false
	}
	tradableMin := now.Hour()*60 + now.Minute()
	return tradableMin >= 9*60+1 && tradableMin <= 15*60+10
}

func (r *Runner) IsRunning() bool {
	return r.running.Load()
}

func (r *Runner) IsEnabled() bool {
	return r.enabled.Load()
}
