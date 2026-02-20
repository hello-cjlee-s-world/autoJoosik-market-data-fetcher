package strategyengine

import "math"

type TechnicalFactor struct{}

func (TechnicalFactor) Name() string { return "technical" }
func (TechnicalFactor) Evaluate(ctx EvalContext) (float64, string) {
	closes := make([]float64, 0, len(ctx.Candidate.Candles))
	for _, c := range ctx.Candidate.Candles {
		closes = append(closes, c.Close)
	}
	if len(closes) < 60 {
		return 0.2, "insufficient candles"
	}
	ma5, ma20, ma60 := SMA(closes, 5), SMA(closes, 20), SMA(closes, 60)
	rsi := RSI(closes, 14)
	_, _, hist := MACD(closes)
	mid, up, low := Bollinger(closes, 20, 2)
	vwap := VWAP(ctx.Candidate.Candles)

	score := 0.0
	if ma5 > ma20 && ma20 > ma60 {
		score += 0.25
	}
	if rsi >= 45 && rsi <= 70 {
		score += 0.20
	}
	if hist > 0 {
		score += 0.20
	}
	price := closes[len(closes)-1]
	if price >= mid && price <= up && price > low {
		score += 0.15
	}
	if price >= vwap {
		score += 0.20
	}
	return math.Min(1, score), "ma/rsi/macd/bb/atr/vwap"
}

type VolumeFactor struct{}

func (VolumeFactor) Name() string { return "volume" }
func (VolumeFactor) Evaluate(ctx EvalContext) (float64, string) {
	candles := ctx.Candidate.Candles
	if len(candles) < 21 {
		return 0.3, "insufficient volume history"
	}
	latest := candles[len(candles)-1]
	prev := candles[len(candles)-2]
	avg20 := 0.0
	for _, c := range candles[len(candles)-21 : len(candles)-1] {
		avg20 += c.Volume
	}
	avg20 /= 20
	pctVsPrev := 0.0
	if prev.Volume > 0 {
		pctVsPrev = (latest.Volume - prev.Volume) / prev.Volume
	}
	multiple := 0.0
	if avg20 > 0 {
		multiple = latest.Volume / avg20
	}
	priceUp := latest.Close >= prev.Close
	score := 0.0
	if pctVsPrev > 0.3 {
		score += 0.35
	}
	if multiple > 1.8 {
		score += 0.35
	}
	if priceUp && latest.Volume > prev.Volume {
		score += 0.30
	}
	return score, "volume expansion"
}

type FlowFactor struct{}

func (FlowFactor) Name() string { return "flow" }
func (FlowFactor) Evaluate(ctx EvalContext) (float64, string) {
	score := 0.0
	if ctx.Candidate.Flow.ForeignNetBuy > 0 {
		score += 0.5
	}
	if ctx.Candidate.Flow.InstitutionNetBuy > 0 {
		score += 0.5
	}
	return score, "foreign/institution flow"
}

type MarketFactor struct{}

func (MarketFactor) Name() string { return "market" }
func (MarketFactor) Evaluate(ctx EvalContext) (float64, string) {
	m := ctx.Candidate.Market
	if m.PanicDrop {
		return 0, "panic filter"
	}
	if m.KospiChangePct > 0 && m.KosdaqChangePct > 0 {
		return 1, "risk-on"
	}
	if m.KospiChangePct > -0.5 && m.KosdaqChangePct > -0.5 {
		return 0.5, "neutral"
	}
	return 0.2, "weak market"
}

type NewsFactor struct{}

func (NewsFactor) Name() string { return "news" }
func (NewsFactor) Evaluate(ctx EvalContext) (float64, string) {
	if len(ctx.Candidate.News) == 0 {
		return 0.4, "no fresh news"
	}
	acc := 0.0
	for _, n := range ctx.Candidate.News {
		acc += (n.Sentiment + 1) / 2
	}
	return acc / float64(len(ctx.Candidate.News)), "news ai sentiment/event"
}

type FinanceFactor struct{}

func (FinanceFactor) Name() string { return "finance" }
func (FinanceFactor) Evaluate(ctx EvalContext) (float64, string) {
	// 재무데이터는 커넥터 연결 전까지 보수적 중립값 사용.
	return 0.5, "finance stub"
}
