package model

import (
	"autoJoosik-market-data-fetcher/pkg/logger"
	"database/sql"
	"encoding/json"
	"time"
)

// TbTradeUniverse 가상 체결 테이블
type TbTradeUniverseEntity struct {
	StkCd        string         `db:"stk_cd" json:"stk_cd"`
	Market       string         `db:"market" json:"market"`
	Name         string         `db:"name" json:"name"`
	Enabled      bool           `db:"enabled" json:"enabled"`
	Status       string         `db:"status" json:"status"`
	StatusReason sql.NullString `db:"status_reason" json:"status_reason"`
	CreatedAt    time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at" json:"updated_at"`
}

func ToTbTradeUniverseEntity(str string) TbTradeUniverseEntity {
	var entity TbTradeUniverseEntity
	err := json.Unmarshal([]byte(str), &entity)
	if err != nil {
		logger.Error("While doing ToTbTradeUniverseEntity :: ", err.Error())
	}

	return entity
}
