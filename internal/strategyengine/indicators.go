package strategyengine

import "math"

func SMA(values []float64, period int) float64 {
	if period <= 0 || len(values) < period {
		return 0
	}
	sum := 0.0
	for _, v := range values[len(values)-period:] {
		sum += v
	}
	return sum / float64(period)
}

func RSI(closes []float64, period int) float64 {
	if period <= 0 || len(closes) <= period {
		return 50
	}
	gain, loss := 0.0, 0.0
	for i := len(closes) - period; i < len(closes); i++ {
		delta := closes[i] - closes[i-1]
		if delta > 0 {
			gain += delta
		} else {
			loss -= delta
		}
	}
	if loss == 0 {
		return 100
	}
	rs := gain / loss
	return 100 - (100 / (1 + rs))
}

func EMA(values []float64, period int) float64 {
	if len(values) == 0 {
		return 0
	}
	k := 2.0 / (float64(period) + 1)
	ema := values[0]
	for _, v := range values[1:] {
		ema = v*k + ema*(1-k)
	}
	return ema
}

func MACD(closes []float64) (macd, signal, hist float64) {
	if len(closes) < 35 {
		return 0, 0, 0
	}
	macdSeries := make([]float64, 0, len(closes))
	for i := 1; i <= len(closes); i++ {
		fast := EMA(closes[:i], 12)
		slow := EMA(closes[:i], 26)
		macdSeries = append(macdSeries, fast-slow)
	}
	macd = macdSeries[len(macdSeries)-1]
	signal = EMA(macdSeries, 9)
	hist = macd - signal
	return
}

func Bollinger(closes []float64, period int, sigma float64) (mid, upper, lower float64) {
	if len(closes) < period || period <= 0 {
		return 0, 0, 0
	}
	window := closes[len(closes)-period:]
	mid = SMA(closes, period)
	varSum := 0.0
	for _, v := range window {
		varSum += math.Pow(v-mid, 2)
	}
	std := math.Sqrt(varSum / float64(period))
	upper = mid + sigma*std
	lower = mid - sigma*std
	return
}

func ATR(candles []Candle, period int) float64 {
	if len(candles) <= period || period <= 0 {
		return 0
	}
	trs := make([]float64, 0, len(candles)-1)
	for i := 1; i < len(candles); i++ {
		highLow := candles[i].High - candles[i].Low
		highPrev := math.Abs(candles[i].High - candles[i-1].Close)
		lowPrev := math.Abs(candles[i].Low - candles[i-1].Close)
		trs = append(trs, math.Max(highLow, math.Max(highPrev, lowPrev)))
	}
	return SMA(trs, period)
}

func VWAP(candles []Candle) float64 {
	denom, numer := 0.0, 0.0
	for _, c := range candles {
		typical := (c.High + c.Low + c.Close) / 3
		numer += typical * c.Volume
		denom += c.Volume
	}
	if denom == 0 {
		return 0
	}
	return numer / denom
}
