package strategyengine

import "time"

type TradingTimeGate struct{}

func (TradingTimeGate) Name() string { return "trading_time" }
func (TradingTimeGate) Evaluate(ctx EvalContext) (bool, string) {
	now := time.Duration(ctx.Now.Hour())*time.Hour + time.Duration(ctx.Now.Minute())*time.Minute
	if now < mustParseHHMM(ctx.Config.TradingStart) || now > mustParseHHMM(ctx.Config.TradingEnd) {
		return false, "outside trading hours"
	}
	return true, ""
}

type LiquidityGate struct{}

func (LiquidityGate) Name() string { return "liquidity" }
func (LiquidityGate) Evaluate(ctx EvalContext) (bool, string) {
	if ctx.Candidate.Turnover < ctx.Config.MinTurnover {
		return false, "turnover too low"
	}
	if ctx.Candidate.Spread > ctx.Config.MaxSpread {
		return false, "spread too wide"
	}
	return true, ""
}

type CircuitBreakerGate struct{}

func (CircuitBreakerGate) Name() string { return "circuit_breaker" }
func (CircuitBreakerGate) Evaluate(ctx EvalContext) (bool, string) {
	if ctx.Candidate.IsHalted || ctx.Candidate.IsVI {
		return false, "halt/vi"
	}
	return true, ""
}

type CooldownGate struct{}

func (CooldownGate) Name() string { return "cooldown" }
func (CooldownGate) Evaluate(ctx EvalContext) (bool, string) {
	if last, ok := ctx.Account.LastExitByCode[ctx.Candidate.Ticker]; ok {
		if ctx.Now.Sub(last) < time.Duration(ctx.Config.CooldownMinutes)*time.Minute {
			return false, "cooldown"
		}
	}
	return true, ""
}

type PositionLimitGate struct{}

func (PositionLimitGate) Name() string { return "position_limit" }
func (PositionLimitGate) Evaluate(ctx EvalContext) (bool, string) {
	if len(ctx.Account.OpenPositions) >= ctx.Config.MaxConcurrentPositions {
		return false, "max concurrent positions"
	}
	weight := ctx.Candidate.Price / ctx.Account.Equity
	if weight > ctx.Config.MaxPositionWeight {
		return false, "position weight exceeded"
	}
	return true, ""
}
