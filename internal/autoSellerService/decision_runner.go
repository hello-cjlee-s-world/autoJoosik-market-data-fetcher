package autoSellerService

import (
	"autoJoosik-market-data-fetcher/internal/model"
	"autoJoosik-market-data-fetcher/internal/repository"
	"autoJoosik-market-data-fetcher/internal/utils"
	"autoJoosik-market-data-fetcher/pkg/logger"
	"context"
	"strings"
	"time"
)

type watchCandidate struct {
	Candidate model.CandidateEntity
	Score     float64
}

func DecideAndExecute(ctx context.Context, pool repository.DB) error {
	decisionMutex.Lock()
	defer decisionMutex.Unlock()

	cfg := LoadStrategyConfig()
	newsProvider := StubNewsProvider{}
	flowProvider := StubFlowProvider{}
	marketProvider := StubMarketProvider{CrashFilterPct: cfg.Gates.CrashFilterPct}
	market := marketProvider.Snapshot(ctx)

	positions, err := repository.GetHoldingPositions(ctx, pool, 0)
	if err != nil {
		return err
	}

	// Exit/Risk engine
	for _, p := range positions {
		candles, _ := LoadRecentCandles(ctx, pool, p.StkCd, 60)
		ind := BuildIndicators(candles)
		position := Position{StkCd: p.StkCd, Qty: p.Qty, AvgPrice: p.AvgPrice, HighestPrice: p.HighestPrice, BuyTime: p.CreatedAt}
		exit := ShouldSell(position, p.LastPrice, cfg.Risk.FixedStopLossPct, cfg.Risk.TakeProfitPct, cfg.Risk.TrailingStopPct)
		if !exit.Do && p.LastPrice > 0 && ind.ATR14 > 0 {
			atrStopPrice := p.AvgPrice - ind.ATR14*cfg.Risk.ATRStopMultiplier
			if p.LastPrice <= atrStopPrice {
				exit = DecisionResult{Do: true, Reason: "atr_stop"}
			}
		}
		if !exit.Do && isScoreCollapsed(ctx, pool, p.StkCd, cfg.Risk.ScoreCollapseDelta) {
			exit = DecisionResult{Do: true, Reason: "score_collapse"}
		}
		if exit.Do {
			sellQty := decideSellQty(exit.Reason, p.Qty, cfg.Sizing.MaxOrderQty)
			logger.Info("Sell decision", "stkCd", p.StkCd, "reason", exit.Reason, "qty", sellQty)
			if err := Sell(p.StkCd, sellQty); err != nil {
				return err
			}
		}
	}

	candidates, err := repository.GetBuyCandidates(ctx, pool, 0)
	if err != nil {
		return err
	}

	// Watchlist 엔진: 뉴스/거래량 급증/수급
	watchPicks := make([]watchCandidate, 0)
	for _, c := range candidates {
		candles, _ := LoadRecentCandles(ctx, pool, c.StkCd, 60)
		volumeBurst := volumeBurstScore(candles)
		flow := flowProvider.NetBuyScore(ctx, c.StkCd)
		news := newsProvider.SentimentScore(ctx, c.StkCd)
		trend, _ := LoadTradeInfoTrend(ctx, pool, c.StkCd, 20)
		score := news*cfg.Watchlist.NewsWeight + volumeBurst*cfg.Watchlist.VolumeWeight + flow*cfg.Watchlist.FlowWeight + trend.Composite*cfg.Watchlist.TrendWeight
		if score >= cfg.Watchlist.MinScore {
			watchPicks = append(watchPicks, watchCandidate{Candidate: c, Score: score})
		}
	}
	watchPicks = topWatchlist(watchPicks, cfg.Watchlist.MaxPicks)

	// Entry 엔진: Gate + Score 합성
	for _, w := range watchPicks {
		c := w.Candidate
		candles, _ := LoadRecentCandles(ctx, pool, c.StkCd, 120)
		ind := BuildIndicators(candles)
		volumeSignal := volumeBurstScore(candles)
		techScore := technicalScore(ind, c.LastPrice)
		trend, _ := LoadTradeInfoTrend(ctx, pool, c.StkCd, 30)
		aggressiveSetup := isAggressiveShortTermSetup(trend, volumeSignal, techScore, cfg)
		dailyPnL := EstimateDailyPnLPercent(ctx, pool, 0)
		entry := NewEngine()
		entry.AddGate(SimpleGate{name: "trade_time", fn: func(e EvalContext) GateResult {
			if !utils.IsTradableTime(e.Now) {
				return GateResult{Pass: false, Reason: "market_closed"}
			}
			return GateResult{Pass: true}
		}})
		entry.AddGate(SimpleGate{name: "market_crash", fn: func(e EvalContext) GateResult {
			if e.Market.IsCrash {
				return GateResult{Pass: false, Reason: "index_crash_filter"}
			}
			return GateResult{Pass: true}
		}})
		entry.AddGate(SimpleGate{name: "liquidity_spread", fn: func(e EvalContext) GateResult {
			turnover := e.Candidate.LastPrice * lastVolume(candles)
			if turnover < cfg.Gates.MinTurnover {
				return GateResult{Pass: false, Reason: "low_turnover"}
			}
			if e.CurrentSpread > 0 && e.CurrentSpread > cfg.Gates.MaxSpreadBps {
				return GateResult{Pass: false, Reason: "wide_spread"}
			}
			return GateResult{Pass: true}
		}})
		entry.AddGate(SimpleGate{name: "cooldown", fn: func(e EvalContext) GateResult {
			if !canBypassCooldown(e.Now, e.Candidate.LastBuyTime, cfg.CooldownDuration(), cfg.Gates.AggressiveCooldownFactor, aggressiveSetup) {
				return GateResult{Pass: false, Reason: "cooldown"}
			}
			return GateResult{Pass: true}
		}})
		entry.AddGate(SimpleGate{name: "holding_limit", fn: func(e EvalContext) GateResult {
			if e.Candidate.CurrentHoldingCount >= cfg.Gates.MaxHoldingCount {
				return GateResult{Pass: false, Reason: "holding_limit"}
			}
			return GateResult{Pass: true}
		}})
		entry.AddGate(SimpleGate{name: "daily_loss_limit", fn: func(e EvalContext) GateResult {
			if e.DailyPnL <= cfg.Gates.DailyLossLimitPct {
				return GateResult{Pass: false, Reason: "daily_loss_limit"}
			}
			return GateResult{Pass: true}
		}})
		entry.AddGate(SimpleGate{name: "order_duplicate", fn: func(e EvalContext) GateResult {
			if e.RecentOrderOpen {
				return GateResult{Pass: false, Reason: "duplicate_order"}
			}
			return GateResult{Pass: true}
		}})
		entry.AddGate(SimpleGate{name: "vi_halt_filter", fn: func(EvalContext) GateResult {
			return GateResult{Pass: true, Reason: "stub_vi_halt"}
		}})

		entry.AddFactor(WeightedFactor{name: "technical", weight: cfg.Entry.TechnicalWeight, factor: func(e EvalContext) float64 {
			return technicalScore(e.Indicators, e.Candidate.LastPrice)
		}})
		entry.AddFactor(WeightedFactor{name: "volume", weight: cfg.Entry.VolumeWeight, factor: func(e EvalContext) float64 {
			return e.VolumeSignal
		}})
		entry.AddFactor(WeightedFactor{name: "flow", weight: cfg.Entry.FlowWeight, factor: func(e EvalContext) float64 {
			return e.FlowScore
		}})
		entry.AddFactor(WeightedFactor{name: "market", weight: cfg.Entry.MarketWeight, factor: func(e EvalContext) float64 {
			if e.Market.KOSPIChangePct > 0 || e.Market.KOSDAQChangePct > 0 {
				return 1
			}
			return 0.3
		}})
		entry.AddFactor(WeightedFactor{name: "news", weight: cfg.Entry.NewsWeight, factor: func(e EvalContext) float64 {
			return e.NewsScore
		}})
		entry.AddFactor(WeightedFactor{name: "trade_trend", weight: cfg.Entry.TradeTrendWeight, factor: func(EvalContext) float64 {
			return trend.Composite
		}})

		evalCtx := EvalContext{Now: time.Now(), Market: market, Candidate: c, Indicators: ind, NewsScore: newsProvider.SentimentScore(ctx, c.StkCd), FlowScore: flowProvider.NetBuyScore(ctx, c.StkCd), VolumeSignal: volumeBurstScore(candles), RecentOrderOpen: hasOpenOrderRecently(ctx, pool, c.StkCd), DailyPnL: dailyPnL, CurrentSpread: 0}
		evalCtx.VolumeSignal = volumeSignal
		effectiveThreshold := adjustedEntryThreshold(cfg.Entry.ThresholdScore, trend, volumeSignal, cfg.Entry.AggressiveThresholdOffset, cfg.Entry.AggressiveMaxOffset)
		pass, score, reasons := entry.Evaluate(evalCtx)
		if pass && score >= effectiveThreshold {
			buyQty := decideBuyQty(score, effectiveThreshold, cfg.Sizing.MaxOrderQty)
			logger.Info("Buy decision", "stkCd", c.StkCd, "score", score, "qty", buyQty, "reasons", strings.Join(reasons, ","))
			if err := Buy(c.StkCd, buyQty); err != nil {
				return err
			}
			break
		}
		logger.Info("Buy reject", "stkCd", c.StkCd, "score", score, "reasons", strings.Join(reasons, ","))
	}

	return reprocessPendingOrders(ctx, pool)
}

func technicalScore(ind IndicatorSnapshot, price float64) float64 {
	s := 0.0
	if ind.MA5 > ind.MA20 {
		s += 0.2
	}
	if ind.MA20 > ind.MA60 {
		s += 0.2
	}
	if ind.RSI14 >= 45 && ind.RSI14 <= 70 {
		s += 0.2
	}
	if ind.MACD > ind.MACDSignal {
		s += 0.2
	}
	if price >= ind.BBMiddle && price <= ind.BBUpper {
		s += 0.1
	}
	if ind.VWAP > 0 && price >= ind.VWAP {
		s += 0.1
	}
	return clamp(s, 0, 1)
}

func volumeBurstScore(candles []OHLCV) float64 {
	if len(candles) < 2 {
		return 0
	}
	curr := candles[len(candles)-1].Volume
	prev := candles[len(candles)-2].Volume
	avg20 := 0.0
	start := 0
	if len(candles) > 20 {
		start = len(candles) - 20
	}
	for _, c := range candles[start:] {
		avg20 += c.Volume
	}
	avg20 /= float64(len(candles[start:]))
	prevRatio := 0.0
	if prev > 0 {
		prevRatio = curr / prev
	}
	avgRatio := 0.0
	if avg20 > 0 {
		avgRatio = curr / avg20
	}
	return clamp((prevRatio+avgRatio)/4, 0, 1)
}

func hasOpenOrderRecently(ctx context.Context, db repository.DB, stkCd string) bool {
	var count int
	_ = db.QueryRow(ctx, `
SELECT COUNT(*)
FROM tb_virtual_order
WHERE account_id = 0
  AND stk_cd = $1
  AND status IN ('NEW','OPEN','PARTIAL')
  AND created_at >= NOW() - INTERVAL '1 minute'
`, stkCd).Scan(&count)
	return count > 0
}

func isScoreCollapsed(ctx context.Context, db repository.DB, stkCd string, collapseDelta float64) bool {
	var scoreTotal float64
	err := db.QueryRow(ctx, `
		SELECT score_total
		FROM tb_stock_score
		WHERE stk_cd = $1
	`, stkCd).Scan(&scoreTotal)
	if err != nil {
		return false
	}

	// scoreTotal is normalized in [0,1].
	// collapseDelta is configured as a negative value (e.g. -0.35),
	// so a value <= 1 + delta means the score has materially weakened.
	return scoreTotal <= 1+collapseDelta
}

func reprocessPendingOrders(ctx context.Context, db repository.DB) error {
	_, err := db.Exec(ctx, `
UPDATE tb_virtual_order
SET status = 'CANCELED', reason = 'auto-reprocess-timeout', updated_at = NOW()
WHERE status IN ('NEW','OPEN','PARTIAL')
  AND created_at < NOW() - INTERVAL '3 minute'
`)
	return err
}

func lastVolume(candles []OHLCV) float64 {
	if len(candles) == 0 {
		return 0
	}
	return candles[len(candles)-1].Volume
}

func adjustedEntryThreshold(base float64, trend TradeInfoTrend, volumeSignal, offset, maxOffset float64) float64 {
	thresholdOffset := 0.0
	if trend.Composite >= 0.75 {
		thresholdOffset += offset
	}
	if trend.Composite >= 0.85 && volumeSignal >= 0.65 {
		thresholdOffset += offset * 0.5
	}
	if trend.BuyPressure >= 0.62 {
		thresholdOffset += 0.02
	}
	thresholdOffset = clamp(thresholdOffset, 0, maxOffset)
	return clamp(base-thresholdOffset, 0.3, 0.95)
}

func canBypassCooldown(now, lastBuyTime time.Time, cooldown time.Duration, aggressiveFactor float64, aggressiveSetup bool) bool {
	elapsed := now.Sub(lastBuyTime)
	if elapsed >= cooldown {
		return true
	}
	if !aggressiveSetup {
		return false
	}
	minCooldown := time.Duration(float64(cooldown) * clamp(aggressiveFactor, 0.1, 1.0))
	return elapsed >= minCooldown
}

func isAggressiveShortTermSetup(trend TradeInfoTrend, volumeSignal, technical float64, cfg StrategyConfig) bool {
	return trend.Composite >= cfg.Entry.AggressiveTrendFloor &&
		trend.BuyPressure >= cfg.Entry.AggressiveBuyPressureFloor &&
		volumeSignal >= cfg.Entry.AggressiveVolumeFloor &&
		technical >= 0.5
}
