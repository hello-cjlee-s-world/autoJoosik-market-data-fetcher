package autoSellerService

import (
	"testing"
	"time"
)

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
	maxOffset := 0.12
	strongTrend := TradeInfoTrend{Composite: 0.9, BuyPressure: 0.65}

	aggressive := adjustedEntryThreshold(base, strongTrend, 0.7, offset, maxOffset)
	if aggressive != 0.38 {
		t.Fatalf("expected aggressive threshold 0.38, got %.2f", aggressive)
	}

	normal := adjustedEntryThreshold(base, TradeInfoTrend{Composite: 0.7, BuyPressure: 0.55}, 0.4, offset, maxOffset)
	if normal != base {
		t.Fatalf("expected unchanged threshold %.2f, got %.2f", base, normal)
	}
}

func TestCanBypassCooldown_AggressiveSetup(t *testing.T) {
	now := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	lastBuy := now.Add(-4 * time.Minute)
	cooldown := 10 * time.Minute

	if !canBypassCooldown(now, lastBuy, cooldown, 0.35, true) {
		t.Fatalf("expected cooldown bypass for aggressive setup")
	}

	if canBypassCooldown(now, lastBuy, cooldown, 0.35, false) {
		t.Fatalf("expected cooldown not to bypass without aggressive setup")
	}
}
