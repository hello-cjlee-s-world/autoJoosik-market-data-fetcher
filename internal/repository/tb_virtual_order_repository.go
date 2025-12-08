package repository

import (
	"autoJoosik-market-data-fetcher/internal/model"
	_interface "autoJoosik-market-data-fetcher/internal/repository/interface"
	"autoJoosik-market-data-fetcher/pkg/logger"
	"context"
)

func InsertOrder(ctx context.Context, db _interface.DB, entity model.TbVirtualOrder) (int64, error) {
	var orderID int64
	err := db.QueryRow(ctx, `
	INSERT INTO tb_virtual_order (
	 user_id, account_id, stk_cd, market, side, order_type, 
	 time_in_force, price, qty, filled_qty, remaining_qty, status, 
	 client_order_id, reason, created_at, updated_at           
	) VALUES (
	   $1, $2, $3, $4, $5,
        $6, $7, $8, $9,
        $10, $11, $12, $13, $14, NOW(), NOW()   
	)
	RETURNING order_id -- order_id 리턴을 위한 코드
-- 	ON CONFLICT (user_id, account_id, cntr_trde_qty) DO NOTHING
`,
		entity.UserID,
		entity.AccountID,
		entity.StkCd,
		entity.Market,
		entity.Side,
		entity.OrderType,
		entity.TimeInForce,
		entity.Price,
		entity.Qty,
		entity.FilledQty,
		entity.RemainingQty,
		entity.Status,
		entity.ClientOrderID,
		entity.Reason,
	).Scan(&orderID)

	if err != nil {
		logger.Error("InsertOrder :: error :: ", err)
		return -1, err
	}
	logger.Debug("InsertOrder :: success :: ", "stk_cd", entity.StkCd)
	return orderID, nil
}
