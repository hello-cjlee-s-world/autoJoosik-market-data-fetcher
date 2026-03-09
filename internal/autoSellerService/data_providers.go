package autoSellerService

import (
	"autoJoosik-market-data-fetcher/internal/repository"
	"context"
	"math"
	"time"
)

type TradeInfoTrend struct {
	PriceMomentum      float64
	StrengthMomentum   float64
	BuyPressure        float64
	VolumeAcceleration float64
	Composite          float64
}

type tradeTrendRow struct {
	price    float64
	strength float64
	buyBid   float64
	selBid   float64
	qty      float64
}

type NewsProvider interface {
	SentimentScore(ctx context.Context, stkCd string) float64
}

type FlowProvider interface {
	NetBuyScore(ctx context.Context, stkCd string) float64
}

type MarketProvider interface {
	Snapshot(ctx context.Context) MarketSnapshot
}

type StubNewsProvider struct{}

func (StubNewsProvider) SentimentScore(context.Context, string) float64 { return 0.5 }

type StubFlowProvider struct{}

func (StubFlowProvider) NetBuyScore(context.Context, string) float64 { return 0.5 }

type StubMarketProvider struct{ CrashFilterPct float64 }

func (s StubMarketProvider) Snapshot(context.Context) MarketSnapshot {
	return MarketSnapshot{KOSPIChangePct: 0, KOSDAQChangePct: 0, IsCrash: false}
}

func LoadRecentCandles(ctx context.Context, db repository.DB, stkCd string, limit int) ([]OHLCV, error) {
	rows, err := db.Query(ctx, `
SELECT cur_prc, cntr_trde_qty
FROM tb_trade_info_log
WHERE stk_cd = $1
ORDER BY tm DESC
LIMIT $2
`, stkCd, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]OHLCV, 0, limit)
	for rows.Next() {
		var closeP, vol float64
		if err := rows.Scan(&closeP, &vol); err != nil {
			return nil, err
		}
		out = append(out, OHLCV{Close: closeP, High: closeP, Low: closeP, Volume: vol})
	}
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	if len(out) == 0 {
		out = append(out, OHLCV{Close: 0, High: 0, Low: 0, Volume: 0})
	}
	return out, rows.Err()
}

func EstimateDailyPnLPercent(ctx context.Context, db repository.DB, accountID int64) float64 {
	var startBalance, currentBalance float64
	_ = db.QueryRow(ctx, `SELECT COALESCE(cash_balance,0) FROM tb_virtual_account WHERE account_id = $1`, accountID).Scan(&currentBalance)
	startOfDay := time.Now().Truncate(24 * time.Hour)
	_ = db.QueryRow(ctx, `
SELECT COALESCE(total_assets, cash_balance, 0)
FROM tb_virtual_asset_daily
WHERE created_at >= $1
ORDER BY created_at ASC
LIMIT 1
`, startOfDay).Scan(&startBalance)
	if startBalance <= 0 {
		return 0
	}
	return (currentBalance - startBalance) / startBalance * 100
}

func LoadTradeInfoTrend(ctx context.Context, db repository.DB, stkCd string, limit int) (TradeInfoTrend, error) {
	if limit < 12 {
		limit = 12
	}
	rows, err := db.Query(ctx, `
SELECT cur_prc, cntr_str, pri_buy_bid_unit, pri_sel_bid_unit, cntr_trde_qty
FROM tb_trade_info_log
WHERE stk_cd = $1
ORDER BY tm DESC
LIMIT $2
`, stkCd, limit)
	if err != nil {
		return TradeInfoTrend{}, err
	}
	defer rows.Close()

	series := make([]tradeTrendRow, 0, limit)
	for rows.Next() {
		var r tradeTrendRow
		if err := rows.Scan(&r.price, &r.strength, &r.buyBid, &r.selBid, &r.qty); err != nil {
			return TradeInfoTrend{}, err
		}
		series = append(series, r)
	}
	if err := rows.Err(); err != nil {
		return TradeInfoTrend{}, err
	}
	for i, j := 0, len(series)-1; i < j; i, j = i+1, j-1 {
		series[i], series[j] = series[j], series[i]
	}

	return summarizeTradeTrend(series), nil
}

func summarizeTradeTrend(series []tradeTrendRow) TradeInfoTrend {
	if len(series) < 2 {
		return TradeInfoTrend{}
	}

	head := series[:len(series)/2]
	tail := series[len(series)/2:]
	first := series[0].price
	last := series[len(series)-1].price
	priceMomentum := 0.0
	if first > 0 {
		priceMomentum = (last - first) / first
	}
	priceScore := clamp((clamp(priceMomentum, -0.02, 0.02)+0.02)/0.04, 0, 1)

	headStrength := 0.0
	tailStrength := 0.0
	for _, r := range head {
		headStrength += r.strength
	}
	for _, r := range tail {
		tailStrength += r.strength
	}
	headStrength /= math.Max(float64(len(head)), 1)
	tailStrength /= math.Max(float64(len(tail)), 1)
	strengthDelta := tailStrength - headStrength
	strengthScore := clamp((clamp(strengthDelta, -20, 20)+20)/40, 0, 1)

	buyPressure := 0.5
	buyTotal := 0.0
	sellTotal := 0.0
	for _, r := range series {
		buyTotal += math.Max(r.buyBid, 0)
		sellTotal += math.Max(r.selBid, 0)
	}
	if buyTotal+sellTotal > 0 {
		buyPressure = buyTotal / (buyTotal + sellTotal)
	}

	headQty := 0.0
	tailQty := 0.0
	for _, r := range head {
		headQty += r.qty
	}
	for _, r := range tail {
		tailQty += r.qty
	}
	headQty /= math.Max(float64(len(head)), 1)
	tailQty /= math.Max(float64(len(tail)), 1)
	volRatio := 1.0
	if headQty > 0 {
		volRatio = tailQty / headQty
	}
	volumeScore := clamp((clamp(volRatio, 0.5, 2.0)-0.5)/1.5, 0, 1)

	composite := clamp(priceScore*0.35+strengthScore*0.25+buyPressure*0.2+volumeScore*0.2, 0, 1)
	return TradeInfoTrend{
		PriceMomentum:      priceScore,
		StrengthMomentum:   strengthScore,
		BuyPressure:        buyPressure,
		VolumeAcceleration: volumeScore,
		Composite:          composite,
	}
}
