package strategyengine

import "time"

type RiskEngine struct{}

func (RiskEngine) Evaluate(ctx EvalContext, position Position, entryScore float64) ExitDecision {
	decision := ExitDecision{}
	price := position.CurrentPrice
	if price <= 0 {
		price = ctx.Candidate.Price
	}

	if ctx.Account.Equity > 0 {
		if -ctx.Account.DailyPnL/ctx.Account.Equity >= ctx.Config.DailyLossLimitPct {
			decision.ExitAll = true
			decision.Reason = append(decision.Reason, "daily loss limit")
		}
	}

	if price <= position.EntryPrice*(1-ctx.Config.FixedStopLossPct) {
		decision.ExitAll = true
		decision.Reason = append(decision.Reason, "fixed stop loss")
	}

	atr := ATR(ctx.Candidate.Candles, 14)
	if atr > 0 && price <= position.EntryPrice-atr*ctx.Config.ATRStopMultiplier {
		decision.ExitAll = true
		decision.Reason = append(decision.Reason, "atr stop loss")
	}

	if !position.TakeProfitDone && price >= position.EntryPrice*(1+ctx.Config.TakeProfitPct) {
		decision.ExitPartial = true
		decision.Reason = append(decision.Reason, "take profit partial")
	}

	if position.PeakPrice > 0 && price <= position.PeakPrice*(1-ctx.Config.TrailingStopPct) {
		decision.ExitAll = true
		decision.Reason = append(decision.Reason, "trailing stop")
	}

	if entryScore < ctx.Config.ExitScoreCollapse {
		decision.ExitAll = true
		decision.Reason = append(decision.Reason, "score collapse")
	}

	if last, ok := ctx.Account.PendingOrders[position.Ticker]; ok {
		if ctx.Now.Sub(last) < time.Duration(ctx.Config.OrderRetrySeconds)*time.Second {
			decision.Reason = append(decision.Reason, "duplicate order blocked")
		} else {
			decision.RetryOrder = true
			decision.Reason = append(decision.Reason, "retry stale unfilled order")
		}
	}

	return decision
}
