package strategyengine

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

type EngineConfig struct {
	EntryThreshold         float64
	ExitScoreCollapse      float64
	Weights                map[string]float64
	RequiredGates          []string
	TradingStart           string
	TradingEnd             string
	MinTurnover            float64
	MaxSpread              float64
	CooldownMinutes        int
	MaxConcurrentPositions int
	MaxPositionWeight      float64
	FixedStopLossPct       float64
	ATRStopMultiplier      float64
	TakeProfitPct          float64
	TrailingStopPct        float64
	DailyLossLimitPct      float64
	OrderRetrySeconds      int
}

func DefaultEngineConfig() EngineConfig {
	return EngineConfig{
		EntryThreshold:         0.55,
		ExitScoreCollapse:      0.30,
		Weights:                map[string]float64{"technical": 0.30, "volume": 0.20, "flow": 0.15, "market": 0.10, "news": 0.15, "finance": 0.10},
		RequiredGates:          []string{"trading_time", "liquidity", "circuit_breaker", "cooldown", "position_limit"},
		TradingStart:           "09:00",
		TradingEnd:             "15:20",
		MinTurnover:            300000000,
		MaxSpread:              0.005,
		CooldownMinutes:        15,
		MaxConcurrentPositions: 8,
		MaxPositionWeight:      0.20,
		FixedStopLossPct:       0.03,
		ATRStopMultiplier:      2.0,
		TakeProfitPct:          0.06,
		TrailingStopPct:        0.025,
		DailyLossLimitPct:      0.03,
		OrderRetrySeconds:      5,
	}
}

// LoadEngineConfig는 strategyEngine 섹션의 핵심 key만 경량 파싱한다.
func LoadEngineConfig(path string) (EngineConfig, error) {
	cfg := DefaultEngineConfig()
	f, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer f.Close()

	inStrategy, inWeights, inGates := false, false, false
	cfg.RequiredGates = nil
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "strategyEngine:") {
			inStrategy = true
			continue
		}
		if !inStrategy {
			continue
		}
		if !strings.HasPrefix(s.Text(), "  ") {
			break
		}
		if strings.HasPrefix(line, "weights:") {
			inWeights, inGates = true, false
			continue
		}
		if strings.HasPrefix(line, "requiredGates:") {
			inGates, inWeights = true, false
			continue
		}
		if inWeights && strings.HasPrefix(s.Text(), "    ") {
			parts := strings.SplitN(strings.TrimSpace(line), ":", 2)
			if len(parts) == 2 {
				if v, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64); err == nil {
					cfg.Weights[parts[0]] = v
				}
			}
			continue
		}
		if inGates && strings.HasPrefix(strings.TrimSpace(line), "-") {
			cfg.RequiredGates = append(cfg.RequiredGates, strings.TrimSpace(strings.TrimPrefix(line, "-")))
			continue
		}
		inWeights, inGates = false, false
		kv := strings.SplitN(line, ":", 2)
		if len(kv) != 2 {
			continue
		}
		k := strings.TrimSpace(kv[0])
		v := strings.Trim(strings.TrimSpace(kv[1]), "\"")
		switch k {
		case "entryThreshold":
			cfg.EntryThreshold, _ = strconv.ParseFloat(v, 64)
		case "exitScoreCollapse":
			cfg.ExitScoreCollapse, _ = strconv.ParseFloat(v, 64)
		case "tradingStart":
			cfg.TradingStart = v
		case "tradingEnd":
			cfg.TradingEnd = v
		case "minTurnover":
			cfg.MinTurnover, _ = strconv.ParseFloat(v, 64)
		case "maxSpread":
			cfg.MaxSpread, _ = strconv.ParseFloat(v, 64)
		case "cooldownMinutes":
			i, _ := strconv.Atoi(v)
			cfg.CooldownMinutes = i
		case "maxConcurrentPositions":
			i, _ := strconv.Atoi(v)
			cfg.MaxConcurrentPositions = i
		case "maxPositionWeight":
			cfg.MaxPositionWeight, _ = strconv.ParseFloat(v, 64)
		case "fixedStopLossPct":
			cfg.FixedStopLossPct, _ = strconv.ParseFloat(v, 64)
		case "atrStopMultiplier":
			cfg.ATRStopMultiplier, _ = strconv.ParseFloat(v, 64)
		case "takeProfitPct":
			cfg.TakeProfitPct, _ = strconv.ParseFloat(v, 64)
		case "trailingStopPct":
			cfg.TrailingStopPct, _ = strconv.ParseFloat(v, 64)
		case "dailyLossLimitPct":
			cfg.DailyLossLimitPct, _ = strconv.ParseFloat(v, 64)
		case "orderRetrySeconds":
			i, _ := strconv.Atoi(v)
			cfg.OrderRetrySeconds = i
		}
	}
	if len(cfg.RequiredGates) == 0 {
		cfg.RequiredGates = []string{"trading_time", "liquidity", "circuit_breaker", "cooldown", "position_limit"}
	}
	return cfg, s.Err()
}

func mustParseHHMM(hhmm string) time.Duration {
	t, _ := time.Parse("15:04", hhmm)
	return time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute
}
