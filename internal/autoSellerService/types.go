package autoSellerService

import "time"

type MarketState struct {
	IsBull      bool // 상승장 판단
	IsBear      bool
	IndexChange float64
	Reason      string
	Volatility  float64 // 변동성
	IsEmergency bool    // 급변 여부
}

type StockState struct {
	StkCd        string
	CurrentPrice float64
	Score        float64 // 종합 점수
}

type Position struct {
	StkCd        string
	Qty          float64
	AvgPrice     float64
	HighestPrice float64
	BuyTime      time.Time
}

type BuyConstraints struct {
	MaxHoldingCount      int
	MaxDailyBuyCount     int
	CooldownAfterBuy     time.Duration
	AllowAddBuy          bool
	MaxInvestPerStockPct float64
}

type DecisionResult struct {
	Do     bool
	Reason string
}
