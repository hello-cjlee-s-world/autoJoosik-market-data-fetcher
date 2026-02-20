package strategyengine

import "fmt"

type TradingEngine struct {
	watchlist *WatchlistEngine
	registry  *EvaluatorRegistry
	risk      RiskEngine
	config    EngineConfig
}

func NewTradingEngine(config EngineConfig, watchlist *WatchlistEngine) *TradingEngine {
	registry := NewEvaluatorRegistry()
	registry.RegisterGate(TradingTimeGate{})
	registry.RegisterGate(LiquidityGate{})
	registry.RegisterGate(CircuitBreakerGate{})
	registry.RegisterGate(CooldownGate{})
	registry.RegisterGate(PositionLimitGate{})
	registry.RegisterFactor(TechnicalFactor{})
	registry.RegisterFactor(VolumeFactor{})
	registry.RegisterFactor(FlowFactor{})
	registry.RegisterFactor(MarketFactor{})
	registry.RegisterFactor(NewsFactor{})
	registry.RegisterFactor(FinanceFactor{})
	return &TradingEngine{watchlist: watchlist, registry: registry, risk: RiskEngine{}, config: config}
}

func (e *TradingEngine) EvaluateEntry(nowCtx EvalContext) EntryDecision {
	nowCtx.Config = e.config
	return e.registry.EvaluateEntry(nowCtx)
}

func (e *TradingEngine) EvaluateExit(nowCtx EvalContext, position Position, latestEntryScore float64) ExitDecision {
	nowCtx.Config = e.config
	return e.risk.Evaluate(nowCtx, position, latestEntryScore)
}

func (e *TradingEngine) Run(rawNews []string, pool []Candidate, account AccountState) map[string]string {
	out := map[string]string{}
	watchlist := e.watchlist.Build(rawNews, pool)
	for _, item := range watchlist {
		candidate, ok := findCandidate(pool, item.Ticker)
		if !ok {
			continue
		}
		entry := e.EvaluateEntry(EvalContext{Now: candidate.Candles[len(candidate.Candles)-1].Time, Candidate: candidate, Account: account})
		if entry.Allow {
			out[item.Ticker] = fmt.Sprintf("ENTRY score=%.2f", entry.Score)
			continue
		}
		out[item.Ticker] = fmt.Sprintf("SKIP score=%.2f reasons=%v", entry.Score, entry.Reason)
	}
	for ticker, pos := range account.OpenPositions {
		candidate, ok := findCandidate(pool, ticker)
		if !ok {
			continue
		}
		exit := e.EvaluateExit(EvalContext{Now: candidate.Candles[len(candidate.Candles)-1].Time, Candidate: candidate, Account: account}, pos, 0.2)
		if exit.ExitAll || exit.ExitPartial {
			out[ticker] = fmt.Sprintf("EXIT all=%v partial=%v reasons=%v", exit.ExitAll, exit.ExitPartial, exit.Reason)
		}
	}
	return out
}

func findCandidate(candidates []Candidate, ticker string) (Candidate, bool) {
	for _, c := range candidates {
		if c.Ticker == ticker {
			return c, true
		}
	}
	return Candidate{}, false
}
