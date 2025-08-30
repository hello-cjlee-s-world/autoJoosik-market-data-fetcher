package database

import (
	"database/sql"
	"fmt"
	"log"

	"autoJoosik-market-data-fetcher/pkg/properties"
	_ "github.com/lib/pq"
)

func Connect() *sql.DB {
	props := properties.GetInstance()

	dbConf := props.Database
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbConf.Host, dbConf.Port, dbConf.User, dbConf.Password, dbConf.DBName, dbConf.SSLMode,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("PostgreSQL 핸들러 생성 실패:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("PostgreSQL 연결 실패:", err)
	}

	log.Println("PostgreSQL 연결 성공!")
	return db
}
