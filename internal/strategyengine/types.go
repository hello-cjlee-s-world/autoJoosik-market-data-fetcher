package strategyengine

import "time"

type Candle struct {
	Time   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}

type FlowSnapshot struct {
	ForeignNetBuy     float64
	InstitutionNetBuy float64
}

type MarketSnapshot struct {
	KospiChangePct  float64
	KosdaqChangePct float64
	PanicDrop       bool
}

type NewsSignal struct {
	Ticker        string
	Sentiment     float64
	EventType     string
	Headline      string
	PublishedTime time.Time
}

type Candidate struct {
	Ticker   string
	Price    float64
	Candles  []Candle
	Flow     FlowSnapshot
	Market   MarketSnapshot
	Spread   float64
	Turnover float64
	News     []NewsSignal
	IsHalted bool
	IsVI     bool
}

type Position struct {
	Ticker         string
	EntryPrice     float64
	CurrentPrice   float64
	Quantity       int
	PeakPrice      float64
	EntryTime      time.Time
	StopLossPrice  float64
	TakeProfitDone bool
	LastOrderAt    time.Time
}

type AccountState struct {
	Equity         float64
	DailyPnL       float64
	OpenPositions  map[string]Position
	PendingOrders  map[string]time.Time
	LastExitByCode map[string]time.Time
}

type EvalContext struct {
	Now            time.Time
	Candidate      Candidate
	Account        AccountState
	Config         EngineConfig
	RunningSession string
}

type WatchlistItem struct {
	Ticker string
	Reason []string
	Score  float64
}

type EntryDecision struct {
	Allow  bool
	Score  float64
	Reason []string
}

type ExitDecision struct {
	ExitAll     bool
	ExitPartial bool
	Reason      []string
	RetryOrder  bool
}
