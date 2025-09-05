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
	"sync"
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

	// db 연결
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

	rst, err := kiwoomApi.GetStockInfo()
	if err == nil {
		err = repository.UpsertStockInfo(context.Background(), datasource.GetPool(), model.ToStockInfoEntity(rst))
	}

	stkCd := "005930"
	rst, err = kiwoomApi.GetTradeInfoLog(stkCd)
	if err == nil {
		fmt.Println("GetTradeInfoLog", rst)
		err = repository.UpsertTradeInfoBatch(context.Background(), datasource.GetPool(), model.ToTradeInfoLogEntity(rst, stkCd))
		if err != nil {
			fmt.Println("UpsertTradeInfoBatch", err.Error())
		}
	}

	TaskConfig := scheduler.LoadTaskConfigs(context.Background(), datasource.GetPool())
	fmt.Println(TaskConfig)
}

func getSchedule(ctx context.Context, pool pgxpool.Pool) {
	//s := gocron.NewScheduler(time.Local)

	// 동적으로 로드
	//taskConfigs := scheduler.LoadTaskConfigs(ctx, pool)
	//for _, cfg := range taskConfigs {
	//	switch {
	//	// every N seconds/minutes 등 처리
	//	case strings.HasPrefix(cfg.Schedule, "every "):
	//		parts := strings.Split(cfg.Schedule, " ")
	//		if len(parts) != 2 {
	//			fmt.Println("잘못된 스케줄:", cfg.Schedule)
	//			continue
	//		}
	//		if parts[1] == "10s" {
	//			s.Every(10).Seconds().Do(getTaskFunc(cfg.TaskType))
	//		}
	//
	//	// 특정 시간 실행 (ex: "09:00")
	//	case strings.Contains(cfg.Schedule, ":"):
	//		s.Every(1).Day().At(cfg.Schedule).Do(getTaskFunc(cfg.TaskType))
	//
	//	default:
	//		fmt.Println("지원하지 않는 스케줄:", cfg.Schedule)
	//	}
	//}
	//
	//s.StartAsync()
}
