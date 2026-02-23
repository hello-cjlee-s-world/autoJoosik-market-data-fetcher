package autoSellerService

import (
	"autoJoosik-market-data-fetcher/internal/model"
	"autoJoosik-market-data-fetcher/pkg/logger"
	"math"
	"sort"
	"strings"
	"time"
)

type Gate interface {
	Name() string
	Evaluate(ctx EvalContext) GateResult
}

type Factor interface {
	Name() string
	Evaluate(ctx EvalContext) FactorResult
}

type GateResult struct {
	Pass   bool
	Reason string
}

type FactorResult struct {
	Score  float64
	Reason string
}

type EvalContext struct {
	Now             time.Time
	Market          MarketSnapshot
	Candidate       model.CandidateEntity
	Position        *Position
	Indicators      IndicatorSnapshot
	NewsScore       float64
	FlowScore       float64
	VolumeSignal    float64
	RecentOrderOpen bool
	DailyPnL        float64
	CurrentSpread   float64
}

type MarketSnapshot struct {
	KOSPIChangePct  float64
	KOSDAQChangePct float64
	IsCrash         bool
}

type Engine struct {
	gates   []Gate
	factors []Factor
}

func NewEngine() *Engine {
	return &Engine{gates: make([]Gate, 0), factors: make([]Factor, 0)}
}

func (e *Engine) AddGate(g Gate)     { e.gates = append(e.gates, g) }
func (e *Engine) AddFactor(f Factor) { e.factors = append(e.factors, f) }

func (e *Engine) Evaluate(ctx EvalContext) (bool, float64, []string) {
	reasons := make([]string, 0)
	for _, g := range e.gates {
		out := g.Evaluate(ctx)
		if !out.Pass {
			reasons = append(reasons, "gate:"+g.Name()+":"+out.Reason)
			return false, 0, reasons
		}
	}

	total := 0.0
	for _, f := range e.factors {
		out := f.Evaluate(ctx)
		total += out.Score
		if out.Reason != "" {
			reasons = append(reasons, "factor:"+f.Name()+":"+out.Reason)
		}
	}
	return true, total, reasons
}

type WeightedFactor struct {
	factor   func(EvalContext) float64
	name     string
	weight   float64
	minScore float64
}

func (w WeightedFactor) Name() string { return w.name }
func (w WeightedFactor) Evaluate(ctx EvalContext) FactorResult {
	raw := w.factor(ctx)
	s := raw * w.weight
	if s < w.minScore {
		return FactorResult{Score: s, Reason: "below_reference"}
	}
	return FactorResult{Score: s}
}

type SimpleGate struct {
	name string
	fn   func(EvalContext) GateResult
}

func (g SimpleGate) Name() string                        { return g.name }
func (g SimpleGate) Evaluate(ctx EvalContext) GateResult { return g.fn(ctx) }

func clamp(v, low, high float64) float64 {
	return math.Max(low, math.Min(high, v))
}

func topWatchlist(candidates []watchCandidate, max int) []watchCandidate {
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].Score == candidates[j].Score {
			return strings.Compare(candidates[i].Candidate.StkCd, candidates[j].Candidate.StkCd) < 0
		}
		return candidates[i].Score > candidates[j].Score
	})
	if len(candidates) > max {
		candidates = candidates[:max]
	}
	for _, c := range candidates {
		logger.Debug("watchlist pick", "stkCd", c.Candidate.StkCd, "score", c.Score)
	}
	return candidates
}
