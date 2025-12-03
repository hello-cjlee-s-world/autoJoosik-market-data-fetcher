package model

import (
	"autoJoosik-market-data-fetcher/pkg/logger"
	"encoding/json"
	"time"
)

// CashLedger 계좌 현금 현황(스냅샷)
type CashLedger struct {
	AcctId       string    `db:"acct_id"`
	SettleCash   float64   `db:"settle_cash"`
	UnsettleCash float64   `db:"unsettle_cash"`
	BuyingPower  float64   `db:"buying_power"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func ToCashLedgerEntity(str string) CashLedger {
	var entity CashLedger
	err := json.Unmarshal([]byte(str), &entity)
	if err != nil {
		logger.Error("While doing ToCashLedgerEntity :: ", err.Error())
	}

	return entity
}
