package autoSellerService

func ShouldSell(
	position Position,
	currentPrice float64,
	stopLossRate float64, // 예: -1.5
	takeProfitRate float64, // 예: +1.0
	trailingRate float64, // 예: -0.8
) DecisionResult {

	if position.Qty <= 0 {
		return DecisionResult{false, "no_position"}
	}

	profitRate := (currentPrice - position.AvgPrice) / position.AvgPrice * 100

	// 손절
	if profitRate <= stopLossRate {
		return DecisionResult{true, "stop_loss"}
	}

	// 트레일링 스탑
	drawDown := (currentPrice - position.HighestPrice) / position.HighestPrice * 100
	if drawDown <= trailingRate {
		return DecisionResult{true, "trailing_stop"}
	}

	// 익절
	if profitRate >= takeProfitRate {
		return DecisionResult{true, "take_profit"}
	}

	return DecisionResult{false, "hold"}
}
