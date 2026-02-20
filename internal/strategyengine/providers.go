package strategyengine

type FlowProvider interface {
	GetFlow(ticker string) FlowSnapshot
}

type MarketProvider interface {
	Snapshot() MarketSnapshot
}

type LiquidityProvider interface {
	GetTurnoverAndSpread(ticker string, price float64) (turnover float64, spread float64)
}

type NewsAIAnalyzer interface {
	Analyze(rawNews []string) []NewsSignal
}

type StubFlowProvider struct{}

func (StubFlowProvider) GetFlow(string) FlowSnapshot {
	return FlowSnapshot{}
}

type StubMarketProvider struct{}

func (StubMarketProvider) Snapshot() MarketSnapshot {
	return MarketSnapshot{}
}

type StubLiquidityProvider struct{}

func (StubLiquidityProvider) GetTurnoverAndSpread(string, float64) (float64, float64) {
	return 0, 0.003
}

type StubNewsAIAnalyzer struct{}

func (StubNewsAIAnalyzer) Analyze(rawNews []string) []NewsSignal {
	res := make([]NewsSignal, 0, len(rawNews))
	for _, n := range rawNews {
		res = append(res, NewsSignal{Headline: n, Sentiment: 0.1, EventType: "general"})
	}
	return res
}
