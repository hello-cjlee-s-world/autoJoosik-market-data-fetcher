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

	//rst, err := kiwoomApi.GetStockInfo()
	//if err == nil {
	//	err = repository.UpsertStockInfo(context.Background(), datasource.GetPool(), model.ToStockInfoEntity(rst))
	//}
	//
	//stkCd := "005930"
	//rst, err = kiwoomApi.GetTradeInfoLog(stkCd)
	//if err == nil {
	//	fmt.Println("GetTradeInfoLog", rst)
	//	err = repository.UpsertTradeInfoBatch(context.Background(), datasource.GetPool(), model.ToTradeInfoLogEntity(rst, stkCd))
	//	if err != nil {
	//		fmt.Println("UpsertTradeInfoBatch", err.Error())
	//	}
	//}

	// scheduler 초기화
	getSchedule(context.Background(), datasource.GetPool())
}

func getSchedule(ctx context.Context, pool *pgxpool.Pool) {
	// 실제 업무 함수들 등록 (task_type -> 함수). 필요 시 pool 캡쳐해서 사용하세요.
	reg := scheduler.Registry{
		"GetTradeInfoLog": func(ctx context.Context) error {
			skdCd := "005930"
			rst, _ := kiwoomApi.GetTradeInfoLog(skdCd)
			ent := model.ToTradeInfoLogEntity(skdCd, rst)
			err := repository.UpsertTradeInfoBatch(ctx, pool, ent)
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
