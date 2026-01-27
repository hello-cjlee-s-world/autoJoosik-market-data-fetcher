package repository

import (
	"autoJoosik-market-data-fetcher/pkg/logger"
	"context"
	"time"
)

func InsertVirtualAssetDaily(ctx context.Context, db DB, total int, stockTotal int, cash int) error {
	var t = time.Now()
	_, err := db.Exec(ctx, `
	INSERT INTO tb_virtual_asset_daily(
		user_id, account_id, base_date, total_assets, stock_value, cash_balance, created_at	                                   
	) values (
		0, 0, $1, $2, $3, $4, now()                                          
	)
    ON CONFLICT (account_id, base_date) DO UPDATE SET
	   total_assets = EXCLUDED.total_assets,
	   stock_value = EXCLUDED.stock_value,
	   cash_balance = EXCLUDED.cash_balance
`, t.Format("2006-01-02"), total, stockTotal, cash)
	if err != nil {
		logger.Error("InsertVirtualAssetDaily :: error :: ", err.Error())
		return err
	}

	logger.Debug("InsertVirtualAssetDaily :: success :: ", "ts", t.Format("2006-01-02 15:04:05"))
	return nil
}
