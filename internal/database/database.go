package database

import (
	"database/sql"
	"fmt"
	"log"

	"autoJoosik-market-data-fetcher/pkg/properties"
	_ "github.com/lib/pq"
)

func DatabaseInit() {
	// database 연결 및 초기화
	db := Connect()
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

// init 하고 config main에서 불러와서 하는걸로 해야함
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
