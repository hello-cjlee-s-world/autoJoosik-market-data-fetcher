package autoSellerService

import "math"

func decideBuyQty(score float64, threshold float64, maxOrderQty int) float64 {
	if maxOrderQty <= 0 {
		maxOrderQty = 100
	}
	if score <= 0 {
		return 1
	}
	if threshold <= 0 {
		threshold = 0.55
	}

	normalized := (score - threshold) / (1 - threshold)
	normalized = clamp(normalized, 0, 1)
	qty := int(math.Round(1 + normalized*float64(maxOrderQty-1)))
	if qty < 1 {
		qty = 1
	}
	if qty > maxOrderQty {
		qty = maxOrderQty
	}
	return float64(qty)
}

func decideSellQty(reason string, positionQty float64, maxOrderQty int) float64 {
	if positionQty <= 0 {
		return 0
	}
	if maxOrderQty <= 0 {
		maxOrderQty = 100
	}

	maxSell := math.Min(positionQty, float64(maxOrderQty))
	if maxSell < 1 {
		return maxSell
	}

	switch reason {
	case "take_profit", "trailing_stop":
		// 익절/트레일링은 절반 청산
		half := math.Floor(maxSell / 2)
		if half < 1 {
			half = 1
		}
		return half
	default:
		// 손절/급락/점수붕괴는 가능한 수량 전량 청산
		return math.Floor(maxSell)
	}
}
