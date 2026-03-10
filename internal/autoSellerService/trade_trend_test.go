package autoSellerService

import "testing"

func TestSummarizeTradeTrend_PositiveMomentum(t *testing.T) {
	series := []tradeTrendRow{
		{price: 100, strength: 90, buyBid: 80, selBid: 120, qty: 1000},
		{price: 101, strength: 95, buyBid: 90, selBid: 110, qty: 1100},
		{price: 103, strength: 110, buyBid: 130, selBid: 90, qty: 1500},
		{price: 104, strength: 120, buyBid: 140, selBid: 80, qty: 1600},
	}

	trend := summarizeTradeTrend(series)
	if trend.Composite <= 0.6 {
		t.Fatalf("expected strong composite trend, got %.4f", trend.Composite)
	}
	if trend.BuyPressure <= 0.5 {
		t.Fatalf("expected buy pressure > 0.5, got %.4f", trend.BuyPressure)
	}
}

func TestAdjustedEntryThreshold_AggressiveOffset(t *testing.T) {
	base := 0.5
	offset := 0.08

	aggressive := adjustedEntryThreshold(base, 0.8, 0.5, IndicatorSnapshot{RSI14: 60, BBMiddle: 100}, 101, offset)
	if aggressive != 0.42 {
		t.Fatalf("expected aggressive threshold 0.42, got %.2f", aggressive)
	}

	normal := adjustedEntryThreshold(base, 0.7, 0.5, IndicatorSnapshot{RSI14: 50, BBMiddle: 100}, 99, offset)
	if normal != base {
		t.Fatalf("expected unchanged threshold %.2f, got %.2f", base, normal)
	}
}

func TestAdjustedEntryThreshold_HighMomentumLowersMore(t *testing.T) {
	base := 0.5
	offset := 0.08
	threshold := adjustedEntryThreshold(base, 0.9, 0.7, IndicatorSnapshot{RSI14: 62, BBMiddle: 100}, 102, offset)
	if threshold != 0.38 {
		t.Fatalf("expected stronger lowered threshold 0.38, got %.2f", threshold)
	}
}

func TestAggressiveBuyQty_BoostedOnStrongMomentum(t *testing.T) {
	qty := aggressiveBuyQty(0.8, 0.5, 100, 0.9, 0.8, IndicatorSnapshot{MA5: 102, MA20: 100, MACD: 1.2, MACDSignal: 1.0, RSI14: 65, BBMiddle: 100, VWAP: 99}, 103)
	if qty <= decideBuyQty(0.8, 0.5, 100) {
		t.Fatalf("expected boosted qty under strong momentum, got %.2f", qty)
	}
}
