package autoSellerService

import (
	"autoJoosik-market-data-fetcher/internal/repository"
	"context"
	"time"
)

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
