package autoSellerService

import (
	"testing"
	"time"
)

func TestShouldBuy_RespectsDurationCooldown(t *testing.T) {
	now := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)

	result := ShouldBuy(
		now,
		MarketState{IsBull: true},
		StockState{Score: 90},
		false,
		now.Add(-5*time.Minute),
		0,
		0,
		BuyConstraints{CooldownAfterBuy: 10 * time.Minute, MaxDailyBuyCount: 10, MaxHoldingCount: 10},
	)

	if result.Do {
		t.Fatalf("expected buy=false due to cooldown, got true")
	}
	if result.Reason != "cool_time" {
		t.Fatalf("expected reason cool_time, got %s", result.Reason)
	}
}

func TestShouldBuy_DeniesAdditionalBuyWhenDisabled(t *testing.T) {
	now := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)

	result := ShouldBuy(
		now,
		MarketState{IsBull: true},
		StockState{Score: 95},
		true,
		now.Add(-30*time.Minute),
		0,
		0,
		BuyConstraints{CooldownAfterBuy: 10 * time.Minute, MaxDailyBuyCount: 10, MaxHoldingCount: 10, AllowAddBuy: false},
	)

	if result.Do {
		t.Fatalf("expected buy=false when additional buy is disabled")
	}
	if result.Reason != "additional_buy_not_allowed" {
		t.Fatalf("expected reason additional_buy_not_allowed, got %s", result.Reason)
	}
}
