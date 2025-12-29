package main

import (
	"autoJoosik-market-data-fetcher/internal/datasource"
	"autoJoosik-market-data-fetcher/internal/kiwoomApi"
	"autoJoosik-market-data-fetcher/internal/scheduler"
	"autoJoosik-market-data-fetcher/pkg/logger"
	"autoJoosik-market-data-fetcher/pkg/properties"
	"context"
	"fmt"
	"github.com/alexflint/go-arg"
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

	// scheduler 초기화
	scheduler.GetSchedule(context.Background(), datasource.GetPool())
}
