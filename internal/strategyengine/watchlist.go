package strategyengine

type WatchlistEngine struct {
	newsAnalyzer NewsAIAnalyzer
	flowProvider FlowProvider
}

func NewWatchlistEngine(newsAnalyzer NewsAIAnalyzer, flowProvider FlowProvider) *WatchlistEngine {
	return &WatchlistEngine{newsAnalyzer: newsAnalyzer, flowProvider: flowProvider}
}

func (w *WatchlistEngine) Build(rawNews []string, candidates []Candidate) []WatchlistItem {
	newsSignals := w.newsAnalyzer.Analyze(rawNews)
	newsByTicker := make(map[string][]NewsSignal)
	for _, n := range newsSignals {
		if n.Ticker == "" {
			continue
		}
		newsByTicker[n.Ticker] = append(newsByTicker[n.Ticker], n)
	}

	items := make([]WatchlistItem, 0)
	for _, c := range candidates {
		flow := w.flowProvider.GetFlow(c.Ticker)
		reasons := make([]string, 0, 3)
		score := 0.0
		if len(c.Candles) > 1 {
			latest := c.Candles[len(c.Candles)-1]
			prev := c.Candles[len(c.Candles)-2]
			if prev.Volume > 0 && latest.Volume/prev.Volume > 1.5 {
				reasons = append(reasons, "volume spike")
				score += 0.4
			}
		}
		if flow.ForeignNetBuy > 0 || flow.InstitutionNetBuy > 0 {
			reasons = append(reasons, "net buying flow")
			score += 0.3
		}
		if _, ok := newsByTicker[c.Ticker]; ok {
			reasons = append(reasons, "news ai trigger")
			score += 0.3
		}
		if len(reasons) > 0 {
			items = append(items, WatchlistItem{Ticker: c.Ticker, Reason: reasons, Score: score})
		}
	}
	return items
}
