package autoSellerService

import (
	"autoJoosik-market-data-fetcher/pkg/properties"
	"time"
)

type StrategyConfig struct {
	Watchlist struct {
		MinScore     float64 `mapstructure:"minScore"`
		MaxPicks     int     `mapstructure:"maxPicks"`
		NewsWeight   float64 `mapstructure:"newsWeight"`
		VolumeWeight float64 `mapstructure:"volumeWeight"`
		FlowWeight   float64 `mapstructure:"flowWeight"`
	} `mapstructure:"watchlist"`
	Entry struct {
		ThresholdScore  float64 `mapstructure:"thresholdScore"`
		TechnicalWeight float64 `mapstructure:"technicalWeight"`
		VolumeWeight    float64 `mapstructure:"volumeWeight"`
		FlowWeight      float64 `mapstructure:"flowWeight"`
		MarketWeight    float64 `mapstructure:"marketWeight"`
		NewsWeight      float64 `mapstructure:"newsWeight"`
	} `mapstructure:"entry"`
	Gates struct {
		MinTurnover       float64 `mapstructure:"minTurnover"`
		MaxSpreadBps      float64 `mapstructure:"maxSpreadBps"`
		CooldownMinutes   int     `mapstructure:"cooldownMinutes"`
		MaxHoldingCount   int     `mapstructure:"maxHoldingCount"`
		MaxPositionPct    float64 `mapstructure:"maxPositionPct"`
		DailyLossLimitPct float64 `mapstructure:"dailyLossLimitPct"`
		CrashFilterPct    float64 `mapstructure:"crashFilterPct"`
	} `mapstructure:"gates"`
	Risk struct {
		FixedStopLossPct   float64 `mapstructure:"fixedStopLossPct"`
		ATRStopMultiplier  float64 `mapstructure:"atrStopMultiplier"`
		TakeProfitPct      float64 `mapstructure:"takeProfitPct"`
		TrailingStopPct    float64 `mapstructure:"trailingStopPct"`
		ScoreCollapseDelta float64 `mapstructure:"scoreCollapseDelta"`
	} `mapstructure:"risk"`
}

func LoadStrategyConfig() StrategyConfig {
	cfg := defaultStrategyConfig()
	_ = properties.UnmarshalKey("strategy", &cfg)
	return cfg
}

func defaultStrategyConfig() StrategyConfig {
	cfg := StrategyConfig{}
	cfg.Watchlist.MinScore = 0.3
	cfg.Watchlist.MaxPicks = 10
	cfg.Watchlist.NewsWeight = 0.2
	cfg.Watchlist.VolumeWeight = 0.5
	cfg.Watchlist.FlowWeight = 0.3

	cfg.Entry.ThresholdScore = 0.55
	cfg.Entry.TechnicalWeight = 0.35
	cfg.Entry.VolumeWeight = 0.2
	cfg.Entry.FlowWeight = 0.2
	cfg.Entry.MarketWeight = 0.15
	cfg.Entry.NewsWeight = 0.1

	cfg.Gates.MinTurnover = 500000000
	cfg.Gates.MaxSpreadBps = 35
	cfg.Gates.CooldownMinutes = 15
	cfg.Gates.MaxHoldingCount = 10
	cfg.Gates.MaxPositionPct = 0.2
	cfg.Gates.DailyLossLimitPct = -3
	cfg.Gates.CrashFilterPct = -2.5

	cfg.Risk.FixedStopLossPct = -1.7
	cfg.Risk.ATRStopMultiplier = 1.2
	cfg.Risk.TakeProfitPct = 2.5
	cfg.Risk.TrailingStopPct = -1.0
	cfg.Risk.ScoreCollapseDelta = -0.35
	return cfg
}

func (s StrategyConfig) CooldownDuration() time.Duration {
	return time.Duration(s.Gates.CooldownMinutes) * time.Minute
}
