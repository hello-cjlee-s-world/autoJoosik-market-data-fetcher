package datasource

import (
	"autoJoosik-market-data-fetcher/pkg/logger"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"time"
)

var DBPool *pgxpool.Pool

type DBConfig struct {
	Url             string
	MaximumPoolSize int
}

var dbConfig DBConfig

func DatasourceInit(conf DBConfig) {
	dbConfig = conf

	var err error
	DBPool, err = pgxpool.NewWithConfig(context.Background(), Config())
	if err != nil {
		logger.Error("Error while creating connection to the database :: ", "error", err.Error())
	}

	// 연결 확인
	if err := DBPool.Ping(context.Background()); err != nil {
		logger.Error("Database ping failed :: ", "error", err.Error())
		panic(err)
	}

	logger.Info("Database connected successfully :: ")
}

func Config() *pgxpool.Config {
	const defaultMaxConns = int32(4)
	const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	poolConfig, err := pgxpool.ParseConfig(dbConfig.Url)
	if err != nil {
		logger.Error("Failed to parse DB config", "error", err.Error())
		panic(err)
	}

	// 풀 사이즈 적용
	if dbConfig.MaximumPoolSize > 0 {
		poolConfig.MaxConns = int32(dbConfig.MaximumPoolSize)
	} else {
		poolConfig.MaxConns = defaultMaxConns
	}

	poolConfig.MinConns = defaultMinConns
	poolConfig.MaxConnLifetime = defaultMaxConnLifetime
	poolConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	poolConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	poolConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	poolConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		logger.Debug("Before acquiring a connection")
		return true
	}

	poolConfig.AfterRelease = func(c *pgx.Conn) bool {
		logger.Debug("After releasing a connection")
		return true
	}

	poolConfig.BeforeClose = func(c *pgx.Conn) {
		logger.Debug("Connection is being closed")
	}

	return poolConfig
}

// 연결 풀 가져오기
func GetPool() *pgxpool.Pool {
	return DBPool
}
