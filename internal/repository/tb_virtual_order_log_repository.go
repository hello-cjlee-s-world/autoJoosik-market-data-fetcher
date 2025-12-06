package repository

import (
	"autoJoosik-market-data-fetcher/internal/model"
	"autoJoosik-market-data-fetcher/pkg/logger"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InsertOrderLog(ctx context.Context, pool *pgxpool.Pool, entity model.TbVirtualOrder) error {
	fmt.Println(entity)
	_, err := pool.Exec(ctx, `
	INSERT INTO tb_virtual_order (
	 user_id, account_id, stk_cd, market, side, order_type, 
	 time_in_force, price, qty, filled_qty, remaining_qty, status, 
	 client_order_id, reason, created_at, updated_at           
	) VALUES (
	   $1, $2, $3, $4, $5,
        $6, $7, $8, $9,
        $10, $11, $12, $13, $14, NOW(), NOW()   
	)
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
	)
	if err != nil {
		logger.Error("insertOrderLog :: error :: ", err)
		return err
	}
	logger.Debug("insertOrderLog :: success :: ", "stk_cd", entity.StkCd)
	return nil
}
