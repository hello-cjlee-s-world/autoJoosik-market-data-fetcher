package model

import "time"

// CandidateEntity 주식 매수 후보
type CandidateEntity struct {
	StkCd               string    `db:"stk_cd" json:"stk_cd"`
	LastPrice           float64   `db:"last_price" json:"last_price"`
	ScoreTotal          float64   `db:"score_total" json:"score_total"`
	AlreadyHolding      bool      `db:"already_holding" json:"already_holding"`
	LastBuyTime         time.Time `db:"last_buy_time" json:"last_buy_time"`
	DailyBuyCount       int       `db:"daily_buy_count" json:"daily_buy_count"`
	CurrentHoldingCount int       `db:"current_holding_count" json:"current_holding_count"`
}
