package main

import (
	"autoJoosik-market-data-fetcher/internal/database"
	"autoJoosik-market-data-fetcher/pkg/logger"
	"autoJoosik-market-data-fetcher/pkg/properties"
	"fmt"
	"github.com/alexflint/go-arg"
	"log"
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
	logger.LoggerInit(logger.Config{
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

	// database 연결 및 초기화
	db := database.Connect()
	defer db.Close()

	// 테스트 쿼리
	rows, err := db.Query("SELECT NOW()")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var now string
		rows.Scan(&now)
		fmt.Println("DB 현재 시간:", now)
	}
}
