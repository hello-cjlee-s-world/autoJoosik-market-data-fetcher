package repository

import (
	"autoJoosik-market-data-fetcher/internal/model"
	"autoJoosik-market-data-fetcher/pkg/logger"
	"context"
)

func InsertTradeLog(ctx context.Context, db DB, entity model.TbVirtualTradeLog) (int64, error) {
	var tradeId int64
	err := db.QueryRow(ctx, `
	INSERT INTO tb_virtual_trade_log (
		order_id, user_id, account_id, stk_cd, market, side, filled_qty, 
	    filled_price, filled_amount, fee_amount, tax_amount, created_at        
	) VALUES (
	   $1, $2, $3, $4, $5,
        $6, $7, $8, $9,
        $10, $11, NOW()   
	)
	RETURNING trade_id -- order_id 리턴을 위한 코드
-- 	ON CONFLICT (user_id, account_id, cntr_trde_qty) DO NOTHING
`,
		entity.OrderID,
		entity.UserID,
		entity.AccountID,
		entity.StkCd,
		entity.Market,
		entity.Side,
		entity.FilledQty,
		entity.FilledPrice,
		entity.FilledAmount,
		entity.FeeAmount,
		entity.TaxAmount,
	).Scan(&tradeId)

	if err != nil {
		logger.Error("InsertTradeLog :: error :: ", err.Error())
		return -1, err
	}
	logger.Debug("InsertTradeLog :: success :: ", "stk_cd", entity.StkCd)
	return tradeId, nil
}
