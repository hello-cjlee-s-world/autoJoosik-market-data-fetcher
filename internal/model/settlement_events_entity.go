package model

import (
	"autoJoosik-market-data-fetcher/pkg/logger"
	"encoding/json"
	"time"
)

// SettlementEventsEntity 정산 이벤트(매도 체결분이 정산되어 현금화되는 스케줄)
type SettlementEventsEntity struct {
	Id          int64     `db:"id"`
	AcctId      string    `db:"acct_id"`
	TradeId     string    `db:"trade_id"`
	StkCd       string    `db:"stk_cd"`
	Side        string    `db:"side"`
	GrossAmount float64   `db:"gross_amount"`
	SettleDate  time.Time `db:"settle_date"`
	CreatedAt   time.Time `db:"created_at"`
}

func ToSettlementEventsEntity(str string) SettlementEventsEntity {
	var entity SettlementEventsEntity
	err := json.Unmarshal([]byte(str), &entity)
	if err != nil {
		logger.Error("While doing ToSettlementEventsEntity :: ", err.Error())
	}

	return entity
}
