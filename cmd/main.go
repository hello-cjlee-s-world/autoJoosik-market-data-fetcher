package main

import (
	"autoJoosik-market-data-fetcher/internal/datasource"
	"autoJoosik-market-data-fetcher/internal/kiwoomApi"
	"autoJoosik-market-data-fetcher/internal/model"
	"autoJoosik-market-data-fetcher/internal/repository"
	"autoJoosik-market-data-fetcher/internal/scheduler"
	"autoJoosik-market-data-fetcher/pkg/logger"
	"autoJoosik-market-data-fetcher/pkg/properties"
	"context"
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type Args struct {
	Config string `arg:"-c, --config" help:"Configuration file"`
}

var props *properties.PropertiesInfo

func main() {
	// 설정 불러오기
	var args Args
	arg.MustParse(&args)
	props = properties.GetInstance()

	if args.Config == "" {
		props.Init("internal/config/autoJoosik_market_data_fetcher_conf.yml")
	} else {
		props.Init(args.Config)
	}

	// logger 초기화
	logger.LoggerInit(logger.LoggerConfig{
		Level:         props.Logging.Level,
		Filename:      props.Logging.Filename,
		MaxSize:       props.Logging.MaxSize,
		MaxBackups:    props.Logging.MaxBackups,
		MaxAge:        props.Logging.MaxAge,
		Compress:      props.Logging.Compress,
		ConsoleOutput: props.Logging.ConsoleOutput,
	})
	logger.Info("server info",
		"port", props.Server.Port,
	)

	var wg sync.WaitGroup
	wg.Add(1)

	// db 연결 초기화
	datasource.DatasourceInit(datasource.DBConfig{
		Url: fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			props.Database.User,
			props.Database.Password,
			props.Database.Host,
			props.Database.Port,
			props.Database.Database,
			props.Database.SSLMode,
		),
		MaximumPoolSize: props.Database.MaximumPoolSize,
	})

	// kiwoom api 초기화
	kiwoomApi.KiwoomInit(kiwoomApi.KiwoomConfig{
		AppKey:    props.KiwoomApi.AppKey,
		SecretKey: props.KiwoomApi.SecretKey,
	})

	// 주식 체결 로직
	stkCd := "005930"
	rst, err := kiwoomApi.GetOrderBookLog(stkCd)
	if err == nil {
		// 트랜잭션으로 묶기
		ctx := context.Background()
		pool := datasource.GetPool()
		tx, err := pool.Begin(ctx)

		orderBookEntity := model.ToOrderBookLogEntity(rst)

		price, _ := strconv.ParseFloat(orderBookEntity.SelFprBid, 64)
		//remainingQty, _ := strconv.ParseFloat(orderBookEntity.SelFprReq, 64)

		// 주문 정보(상태) insert
		qty := 1.0          // 내가 사고 싶은 수량
		remainingQty := 0.0 // 내 주문 기준 남은 수량
		userId := int64(0)
		accountID := int64(0)
		clientOrderID := fmt.Sprintf("%d_%d_%d",
			accountID,
			time.Now().UnixNano(),
			rand.Intn(1000),
		) // 유니크한 주문 아이디
		virtualOrderEntity := model.TbVirtualOrder{
			UserID:       userId,
			AccountID:    accountID,
			StkCd:        stkCd,
			Market:       "KOSPI",
			Side:         "B",
			OrderType:    "MARKET",
			TimeInForce:  "DAY",
			Price:        price,
			Qty:          qty,
			FilledQty:    qty,
			RemainingQty: remainingQty,
			Status:       "FILLED", //'NEW','OPEN','PARTIAL','FILLED','CANCELED','REJECTED'

			ClientOrderID: clientOrderID, // 클라이언트에서 부여하는 주문 ID (선택)
			Reason:        "",            //-- 거절/취소 사유 등 (옵션)
		}
		orderId, err := repository.InsertOrder(ctx, tx, virtualOrderEntity)
		if err != nil {
			fmt.Println("InsertOrderLog", err.Error())
		}

		// 거래 로그 insert(원래는 주문 체결 후 insert 지만 가상이라 바로 체결함에 따라 바로 insert)
		virtualTradeLogEntity := model.TbVirtualTradeLog{
			OrderID:      orderId,
			UserID:       0,
			AccountID:    0,
			StkCd:        stkCd,
			Market:       "KOSPI",
			Side:         "B",
			FilledQty:    qty,
			FilledPrice:  price,
			FilledAmount: qty * price,
			FeeAmount:    0.1,
			TaxAmount:    0.1,
		}

		tradeId, err := repository.InsertTradeLog(ctx, tx, virtualTradeLogEntity)
		if err != nil {
			fmt.Println("InsertTradeLog", err.Error())
		}
		fmt.Println(tradeId)

		// 트랜잭션으로 묶어서 commit
		if err := tx.Commit(ctx); err != nil {

		}

	}

	rst, err = kiwoomApi.GetTradeInfoLog(stkCd)
	if err == nil {
		fmt.Println("GetTradeInfoLog", rst)
		err = repository.UpsertTradeInfoBatch(context.Background(), datasource.GetPool(), model.ToTradeInfoLogEntity(rst, stkCd))
		if err != nil {
			fmt.Println("UpsertTradeInfoBatch", err.Error())
		}
	}

	// scheduler 초기화
	getSchedule(context.Background(), datasource.GetPool())
}

func getSchedule(ctx context.Context, pool *pgxpool.Pool) {
	// 실제 업무 함수들 등록 (task_type -> 함수). 필요 시 pool 캡쳐해서 사용
	reg := scheduler.Registry{
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
	r := scheduler.NewRunner(pool, reg)
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
