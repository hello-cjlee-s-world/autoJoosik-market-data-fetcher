package autoSellerService

import "testing"

func TestDecideBuyQty(t *testing.T) {
	qty := decideBuyQty(0.55, 0.55, 100)
	if qty != 1 {
		t.Fatalf("expected minimum qty 1 at threshold, got %v", qty)
	}

	qty = decideBuyQty(1.0, 0.55, 100)
	if qty != 100 {
		t.Fatalf("expected max qty 100 at max score, got %v", qty)
	}
}

func TestDecideSellQty(t *testing.T) {
	qty := decideSellQty("take_profit", 90, 100)
	if qty != 45 {
		t.Fatalf("expected half qty 45 for take_profit, got %v", qty)
	}

	qty = decideSellQty("stop_loss", 150, 100)
	if qty != 100 {
		t.Fatalf("expected capped full sell qty 100 for stop_loss, got %v", qty)
	}
}
