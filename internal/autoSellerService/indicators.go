package autoSellerService

import "math"

type OHLCV struct {
	High   float64
	Low    float64
	Close  float64
	Volume float64
}

type IndicatorSnapshot struct {
	MA5        float64
	MA20       float64
	MA60       float64
	RSI14      float64
	MACD       float64
	MACDSignal float64
	BBUpper    float64
	BBMiddle   float64
	BBLower    float64
	ATR14      float64
	VWAP       float64
}

func BuildIndicators(candles []OHLCV) IndicatorSnapshot {
	closes := make([]float64, 0, len(candles))
	highs := make([]float64, 0, len(candles))
	lows := make([]float64, 0, len(candles))
	vols := make([]float64, 0, len(candles))
	for _, c := range candles {
		closes = append(closes, c.Close)
		highs = append(highs, c.High)
		lows = append(lows, c.Low)
		vols = append(vols, c.Volume)
	}
	macd, sig := MACD(closes)
	up, mid, low := Bollinger(closes, 20, 2)
	return IndicatorSnapshot{
		MA5:        SMA(closes, 5),
		MA20:       SMA(closes, 20),
		MA60:       SMA(closes, 60),
		RSI14:      RSI(closes, 14),
		MACD:       macd,
		MACDSignal: sig,
		BBUpper:    up,
		BBMiddle:   mid,
		BBLower:    low,
		ATR14:      ATR(highs, lows, closes, 14),
		VWAP:       VWAP(closes, vols),
	}
}

func SMA(values []float64, period int) float64 {
	if len(values) == 0 {
		return 0
	}
	if len(values) < period {
		period = len(values)
	}
	sum := 0.0
	for _, v := range values[len(values)-period:] {
		sum += v
	}
	return sum / float64(period)
}

func EMA(values []float64, period int) float64 {
	if len(values) == 0 {
		return 0
	}
	k := 2.0 / float64(period+1)
	ema := values[0]
	for _, v := range values[1:] {
		ema = v*k + ema*(1-k)
	}
	return ema
}

func RSI(closes []float64, period int) float64 {
	if len(closes) <= period {
		return 50
	}
	gain := 0.0
	loss := 0.0
	for i := len(closes) - period; i < len(closes); i++ {
		d := closes[i] - closes[i-1]
		if d > 0 {
			gain += d
		} else {
			loss -= d
		}
	}
	if loss == 0 {
		return 100
	}
	rs := gain / loss
	return 100 - (100 / (1 + rs))
}

func MACD(closes []float64) (float64, float64) {
	if len(closes) == 0 {
		return 0, 0
	}
	fast := EMA(closes, 12)
	slow := EMA(closes, 26)
	macd := fast - slow
	return macd, EMA([]float64{macd}, 9)
}

func Bollinger(closes []float64, period int, multiplier float64) (float64, float64, float64) {
	if len(closes) == 0 {
		return 0, 0, 0
	}
	if len(closes) < period {
		period = len(closes)
	}
	window := closes[len(closes)-period:]
	mean := SMA(window, len(window))
	var variance float64
	for _, v := range window {
		variance += math.Pow(v-mean, 2)
	}
	std := math.Sqrt(variance / float64(len(window)))
	return mean + std*multiplier, mean, mean - std*multiplier
}

func ATR(highs, lows, closes []float64, period int) float64 {
	if len(closes) < 2 {
		return 0
	}
	trs := make([]float64, 0, len(closes)-1)
	for i := 1; i < len(closes); i++ {
		tr := math.Max(highs[i]-lows[i], math.Max(math.Abs(highs[i]-closes[i-1]), math.Abs(lows[i]-closes[i-1])))
		trs = append(trs, tr)
	}
	return SMA(trs, period)
}

func VWAP(closes, volumes []float64) float64 {
	if len(closes) == 0 || len(closes) != len(volumes) {
		return 0
	}
	num := 0.0
	den := 0.0
	for i := range closes {
		num += closes[i] * volumes[i]
		den += volumes[i]
	}
	if den == 0 {
		return closes[len(closes)-1]
	}
	return num / den
}
