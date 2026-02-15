package autoSellerService

import "testing"

func TestShouldSell_InvalidPriceData(t *testing.T) {
	result := ShouldSell(Position{Qty: 1, AvgPrice: 0, HighestPrice: 0}, 100, -1.5, 1.0, -0.8)

	if result.Do {
		t.Fatalf("expected sell=false for invalid price data")
	}
	if result.Reason != "invalid_price_data" {
		t.Fatalf("expected reason invalid_price_data, got %s", result.Reason)
	}
}

func TestShouldSell_TrailingStop(t *testing.T) {
	result := ShouldSell(Position{Qty: 1, AvgPrice: 100, HighestPrice: 110}, 108, -1.5, 10.0, -0.8)

	if !result.Do {
		t.Fatalf("expected sell=true for trailing stop")
	}
	if result.Reason != "trailing_stop" {
		t.Fatalf("expected reason trailing_stop, got %s", result.Reason)
	}
}
