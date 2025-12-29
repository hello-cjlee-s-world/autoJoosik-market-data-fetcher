// pkg/scheduler/scheduler.go
package scheduler

import (
	"autoJoosik-market-data-fetcher/internal/kiwoomApi"
	"autoJoosik-market-data-fetcher/internal/model"
	"autoJoosik-market-data-fetcher/internal/repository"
	"context"
	"errors"
	"log"
	"strconv"
	"strings"
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
}

// Runner 생성
func NewRunner(pool *pgxpool.Pool, reg Registry) *Runner {
	s := gocron.NewScheduler(time.Local)
	s.TagsUnique()
	return &Runner{pool: pool, registry: reg, s: s}
}

// 시작 (로드 후 비동기 실행)
func (r *Runner) Start(ctx context.Context) error {
	if err := r.Reload(ctx); err != nil {
		return err
	}
	r.s.StartAsync()
	return nil
}

// 중지
func (r *Runner) Stop() {
	r.s.Stop()
}

// 스케줄 전부 리로드 (기존 잡 제거 후 DB에서 다시 설정)
func (r *Runner) Reload(ctx context.Context) error {
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
	rows, err := pool.Query(ctx, `SELECT name, schedule, task_type FROM schedule_info WHERE enabled = true`)
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

func GetSchedule(ctx context.Context, pool *pgxpool.Pool) {
	// 실제 업무 함수들 등록 (task_type -> 함수). 필요 시 pool 캡쳐해서 사용
	reg := Registry{
		"GetTradeInfoLog": func(ctx context.Context) error {
			stkCd := "005930"
			rst, err := kiwoomApi.GetTradeInfoLog(stkCd)
			if err != nil {
				return err
			}
			entList := model.ToTradeInfoLogEntity(rst, stkCd)
			err = repository.UpsertTradeInfoBatch(ctx, pool, entList)
			if err != nil {
				return err
			}
			return nil
		},
		"UpsertStockInfo": func(ctx context.Context) error {
			stkCd := "005930"
			rst, err := kiwoomApi.GetStockInfo(stkCd)
			if err != nil {
				return err
			}
			ent := model.ToStockInfoEntity(rst)
			err = repository.UpsertStockInfo(ctx, pool, ent)
			if err != nil {
				return err
			}
			return nil
		},
		"SellOrBuy": func(ctx context.Context) error {
			stkCd := "005930"
			rst, err := kiwoomApi.GetStockInfo(stkCd)
			if err != nil {
				return err
			}
			ent := model.ToStockInfoEntity(rst)
			err = repository.UpsertStockInfo(ctx, pool, ent)
			if err != nil {
				return err
			}
			return nil
		},
	}

	// 러너 생성 & 시작
	r := NewRunner(pool, reg)
	if err := r.Start(ctx); err != nil {
		log.Fatal("scheduler start:", err)
	}
	log.Println("[scheduler] started")

	// (옵션) 주기적 리로드: schedule_info 변경 반영
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			r.Stop()
			log.Println("[scheduler] stopped")
			return
		case <-ticker.C:
			if err := r.Reload(ctx); err != nil {
				log.Println("[scheduler] reload error:", err)
			}
		}
	}
}
