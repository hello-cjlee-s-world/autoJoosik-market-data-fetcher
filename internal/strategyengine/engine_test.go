package strategyengine

import (
	"testing"
	"time"
)

func TestEvaluateEntryGateAndScore(t *testing.T) {
	cfg := DefaultEngineConfig()
	engine := NewTradingEngine(cfg, NewWatchlistEngine(StubNewsAIAnalyzer{}, StubFlowProvider{}))

	candles := make([]Candle, 70)
	now := time.Date(2025, 1, 2, 10, 0, 0, 0, time.Local)
	for i := range candles {
		candles[i] = Candle{Time: now.Add(time.Duration(i) * time.Minute), Open: 100 + float64(i), High: 101 + float64(i), Low: 99 + float64(i), Close: 100 + float64(i), Volume: 20000 + float64(i*200)}
	}
	candidate := Candidate{
		Ticker:   "005930",
		Price:    candles[len(candles)-1].Close,
		Candles:  candles,
		Flow:     FlowSnapshot{ForeignNetBuy: 1000, InstitutionNetBuy: 800},
		Market:   MarketSnapshot{KospiChangePct: 0.5, KosdaqChangePct: 0.4},
		Turnover: 900000000,
		Spread:   0.002,
	}
	ctx := EvalContext{Now: now, Candidate: candidate, Account: AccountState{Equity: 100000000, OpenPositions: map[string]Position{}, LastExitByCode: map[string]time.Time{}}}
	decision := engine.EvaluateEntry(ctx)
	if !decision.Allow {
		t.Fatalf("expected entry allowed, got %+v", decision)
	}
	if decision.Score < cfg.EntryThreshold {
		t.Fatalf("score below threshold: %.2f", decision.Score)
	}
}

func TestRiskDailyLossExit(t *testing.T) {
	cfg := DefaultEngineConfig()
	r := RiskEngine{}
	now := time.Now()
	ctx := EvalContext{Now: now, Config: cfg, Candidate: Candidate{Price: 9700}, Account: AccountState{Equity: 1000000, DailyPnL: -50000, PendingOrders: map[string]time.Time{}}}
	pos := Position{Ticker: "000660", EntryPrice: 10000, CurrentPrice: 9600, PeakPrice: 10300}
	decision := r.Evaluate(ctx, pos, 0.1)
	if !decision.ExitAll {
		t.Fatalf("expected full exit, got %+v", decision)
	}
}
