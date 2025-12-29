package autoSellerService

import "time"

type MarketState struct {
	IsBull      bool    // 상승장 판단
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
	MaxInvestmentPerStock float64
	MaxDailyBuyCount      int
	MaxConcurrentHolding  int
	CoolTimeMinutes       int
}

type DecisionResult struct {
	Do     bool
	Reason string
}
