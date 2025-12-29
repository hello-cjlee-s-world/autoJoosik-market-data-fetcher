package autoSellerService

import "time"

func ShouldBuy(
	now time.Time,
	market MarketState,
	stock StockState,
	alreadyHolding bool,
	lastBuyTime time.Time,
	dailyBuyCount int,
	currentHoldingCount int,
	constraints BuyConstraints,
) DecisionResult {

	if market.IsEmergency {
		return DecisionResult{false, "market_emergency"}
	}

	if !market.IsBull {
		return DecisionResult{false, "not_bull_market"}
	}

	if stock.Score < 70 {
		return DecisionResult{false, "score_too_low"}
	}

	if alreadyHolding {
		// 추가 매수 조건 (단순 예시)
		if stock.Score < 85 {
			return DecisionResult{false, "additional_buy_condition_not_met"}
		}
	}

	if now.Sub(lastBuyTime) < time.Minute*time.Duration(constraints.CoolTimeMinutes) {
		return DecisionResult{false, "cool_time"}
	}

	if dailyBuyCount >= constraints.MaxDailyBuyCount {
		return DecisionResult{false, "daily_buy_limit"}
	}

	if currentHoldingCount >= constraints.MaxConcurrentHolding {
		return DecisionResult{false, "concurrent_holding_limit"}
	}

	return DecisionResult{true, "buy_signal"}
}
