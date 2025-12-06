package model

import (
	"autoJoosik-market-data-fetcher/pkg/logger"
	"encoding/json"
	"time"
)

// TbVirtualTradeLog 가상 체결 테이블
type TbVirtualTradeLog struct {
	TradeID      int64     `db:"trade_id" json:"trade_id"`
	OrderID      int64     `db:"order_id" json:"order_id"`
	UserID       int64     `db:"user_id" json:"user_id"`
	AccountID    int64     `db:"account_id" json:"account_id"`
	StkCd        string    `db:"stk_cd" json:"stk_cd"`
	Market       string    `db:"market" json:"market"`
	Side         string    `db:"side" json:"side"`
	FilledQty    float64   `db:"filled_qty" json:"filled_qty"`
	FilledPrice  float64   `db:"filled_price" json:"filled_price"`
	FilledAmount float64   `db:"filled_amount" json:"filled_amount"`
	FeeAmount    float64   `db:"fee_amount" json:"fee_amount"`
	TaxAmount    float64   `db:"tax_amount" json:"tax_amount"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

func ToTbVirtualTradeLog(str string) TbVirtualTradeLog {
	var entity TbVirtualTradeLog
	err := json.Unmarshal([]byte(str), &entity)
	if err != nil {
		logger.Error("While doing ToTbVirtualAssetEntity :: ", err.Error())
	}

	return entity
}
