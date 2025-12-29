package model

import "time"

type HoldingPosition struct {
	AccountId int64
	UserId    int64
	StkCd     string
	Market    string

	Qty          float64
	AvailableQty float64

	AvgPrice  float64
	LastPrice float64

	HighestPrice float64 // 매수 이후 최고가 (trailing stop 용)
	InvestedAmt  float64

	CreatedAt time.Time // 최초 매수 시점
	UpdatedAt time.Time // 마지막 갱신 시점
}
