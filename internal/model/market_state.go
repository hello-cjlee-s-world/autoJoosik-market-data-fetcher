package model

type MarketState struct {
	IsBull      bool    // 상승장 여부
	IsBear      bool    // 하락장 여부
	Volatility  float64 // 변동성 지표 (예: ATR, stddev)
	IndexChange float64 // 지수 변화율 (%)
	IsEmergency bool    // 급변/거래중지 여부
	Reason      string  // 판단 사유 (로그용)
}
