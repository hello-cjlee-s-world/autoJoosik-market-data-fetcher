package autoSellerService

import "math"

func buildMarketState(r1, r2, r3, vol float64) MarketState {
	state := MarketState{}

	state.Volatility = vol
	state.IndexChange = r1

	if math.Abs(r1) >= 2.0 || vol >= 2.5 {
		state.IsEmergency = true
		state.Reason = "market_shock"
		return state
	}

	positiveSignals := 0
	if r1 > 0 {
		positiveSignals++
	}
	if r2 > 0 {
		positiveSignals++
	}
	if r3 > 0 {
		positiveSignals++
	}

	trendScore := (r1 * 0.5) + (r2 * 0.3) + (r3 * 0.2)

	// 기존(3개 지표 모두 양수) 조건이 너무 엄격해 매수 기회를 놓치는 경우가 많아
	// 2/3 양수 또는 가중 추세 점수 양수면 상승장으로 완화한다.
	state.IsBull = positiveSignals >= 2 || trendScore > 0
	state.IsBear = r1 < 0 && r2 < 0

	if state.IsBull {
		state.Reason = "bull_relaxed"
	} else {
		state.Reason = "not_bull"
	}

	return state
}
