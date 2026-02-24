package repository

import (
	"autoJoosik-market-data-fetcher/internal/model"
	"autoJoosik-market-data-fetcher/pkg/logger"
	"context"
)

func UpsertVirtualAsset(ctx context.Context, db DB, entity model.TbVirtualAssetEntity) error {
	_, err := db.Exec(ctx, `
		INSERT INTO tb_virtual_asset (
		  user_id, account_id, stk_cd, market, position_side,
		  qty, available_qty, avg_price, last_price, highest_price,
		  invested_amount, eval_amount, eval_pl, eval_pl_rate,
		  today_buy_qty, today_sell_qty,
		  status, last_eval_at,
		  created_at, updated_at
		) VALUES (
		  $1, $2, $3, $4, $7::text,
		  $5, $5, $6, $6, $6,
		  ($5::numeric * $6::numeric),
		  ($5::numeric * $6::numeric),
		  0,
		  0,
		  CASE WHEN $7::text = 'B' THEN $5 ELSE 0 END,
		  CASE WHEN $7::text = 'S' THEN $5 ELSE 0 END,
		  $8, NOW(),
		  NOW(), NOW()
		)
		ON CONFLICT (account_id, stk_cd, market, position_side)
		DO UPDATE SET
		  user_id = EXCLUDED.user_id,

		  qty = CASE
			WHEN $7::text = 'B' THEN tb_virtual_asset.qty + EXCLUDED.qty
			WHEN $7::text = 'S' THEN GREATEST(tb_virtual_asset.qty - EXCLUDED.qty, 0)
			ELSE tb_virtual_asset.qty
		  END,

		  available_qty = CASE
			WHEN $7::text = 'B' THEN tb_virtual_asset.available_qty + EXCLUDED.qty
			WHEN $7::text = 'S' THEN GREATEST(tb_virtual_asset.available_qty - EXCLUDED.qty, 0)
			ELSE tb_virtual_asset.available_qty
		  END,

		  invested_amount = CASE
			WHEN $7::text = 'B' THEN tb_virtual_asset.invested_amount + EXCLUDED.invested_amount
			WHEN $7::text = 'S' THEN GREATEST(tb_virtual_asset.invested_amount - (tb_virtual_asset.avg_price * EXCLUDED.qty), 0)
			ELSE tb_virtual_asset.invested_amount
		  END,

		  avg_price = CASE
			WHEN $7::text = 'B'
			 AND (tb_virtual_asset.qty + EXCLUDED.qty) > 0
			  THEN (tb_virtual_asset.invested_amount + EXCLUDED.invested_amount)
				   / (tb_virtual_asset.qty + EXCLUDED.qty)
			WHEN $7::text = 'S'
			 AND (tb_virtual_asset.qty - EXCLUDED.qty) <= 0 THEN 0
			ELSE tb_virtual_asset.avg_price
		  END,

		  last_price = EXCLUDED.last_price,

		highest_price = CASE
		  WHEN EXCLUDED.position_side = 'B'
		   AND tb_virtual_asset.highest_price < EXCLUDED.avg_price
			THEN EXCLUDED.avg_price
		  WHEN EXCLUDED.position_side = 'S'
			THEN tb_virtual_asset.highest_price
		  ELSE tb_virtual_asset.avg_price
		END,

		eval_amount = CASE
		  WHEN $7::text = 'B' THEN (tb_virtual_asset.qty + EXCLUDED.qty) * EXCLUDED.last_price
		  WHEN $7::text = 'S' THEN GREATEST(tb_virtual_asset.qty - EXCLUDED.qty, 0) * EXCLUDED.last_price
		  ELSE tb_virtual_asset.eval_amount
		END,

		eval_pl = CASE
		  WHEN $7::text = 'B' THEN ((tb_virtual_asset.qty + EXCLUDED.qty) * EXCLUDED.last_price) - (tb_virtual_asset.invested_amount + EXCLUDED.invested_amount)
		  WHEN $7::text = 'S' THEN (GREATEST(tb_virtual_asset.qty - EXCLUDED.qty, 0) * EXCLUDED.last_price)
									  - GREATEST(tb_virtual_asset.invested_amount - (tb_virtual_asset.avg_price * EXCLUDED.qty), 0)
		  ELSE tb_virtual_asset.eval_pl
		END,

		eval_pl_rate = CASE
		  WHEN $7::text = 'B'
		   AND (tb_virtual_asset.invested_amount + EXCLUDED.invested_amount) > 0
			THEN ((((tb_virtual_asset.qty + EXCLUDED.qty) * EXCLUDED.last_price) - (tb_virtual_asset.invested_amount + EXCLUDED.invested_amount))
					 / (tb_virtual_asset.invested_amount + EXCLUDED.invested_amount)) * 100
		  WHEN $7::text = 'S'
		   AND GREATEST(tb_virtual_asset.invested_amount - (tb_virtual_asset.avg_price * EXCLUDED.qty), 0) > 0
			THEN (((GREATEST(tb_virtual_asset.qty - EXCLUDED.qty, 0) * EXCLUDED.last_price)
					 - GREATEST(tb_virtual_asset.invested_amount - (tb_virtual_asset.avg_price * EXCLUDED.qty), 0))
					 / GREATEST(tb_virtual_asset.invested_amount - (tb_virtual_asset.avg_price * EXCLUDED.qty), 0)) * 100
		  ELSE 0
		END,

		  today_buy_qty  = tb_virtual_asset.today_buy_qty  + EXCLUDED.today_buy_qty,
		  today_sell_qty = tb_virtual_asset.today_sell_qty + EXCLUDED.today_sell_qty,

		  status = CASE
			WHEN $7::text = 'S' AND GREATEST(tb_virtual_asset.qty - EXCLUDED.qty, 0) = 0 THEN 'CLOSED'
			ELSE EXCLUDED.status
		  END,
		  last_eval_at = NOW(),
		  updated_at   = NOW();
		`,
		entity.UserId,
		entity.AccountId,
		entity.StkCd,
		entity.Market,
		entity.Qty,
		entity.AvgPrice,
		entity.PositionSide,
		entity.Status,
	)
	if err != nil {
		logger.Error("UpsertVirtualAsset :: error :: ", err.Error())
		return err
	}
	logger.Debug("UpsertVirtualAsset :: success :: ", "stk_cd", entity.StkCd)
	return nil
}

func GetAvailableAssetsByAccountAndStkCd(ctx context.Context, db DB, accountId int64, stkCd string,
) (float64, error) {
	var availableQty float64

	err := db.QueryRow(ctx, `
		SELECT available_qty
		FROM tb_virtual_asset
		WHERE account_id = $1
		  AND stk_cd     = $2
	`, accountId, stkCd).Scan(&availableQty)

	if err != nil {
		logger.Error(
			"GetAvailableAssetsByAccountAndStkCd :: error",
			"accountId", accountId,
			"stkCd", stkCd,
			"err", err.Error(),
		)
		return 0, err
	}
	logger.Debug("GetAvailableAssets... ok", "accountId", accountId, "stkCd", stkCd, "availableQty", availableQty)

	return availableQty, nil
}

func GetHoldingPositions(ctx context.Context, db DB, accountId int64) ([]model.HoldingPositionEntity, error) {
	var holdingPositions []model.HoldingPositionEntity

	rows, err := db.Query(ctx, `
	SELECT 
	    account_id,
		user_id,
		stk_cd,
		market,
		qty,
		available_qty,
		avg_price,
		last_price,
		highest_price,
		invested_amount,
		created_at,
		updated_at
	FROM tb_virtual_asset
	WHERE account_id = $1
`, accountId)

	if err != nil {
		logger.Error("GetBuyCandidates :: error :: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var position model.HoldingPositionEntity
		if err := rows.Scan(
			&position.AccountId,
			&position.UserId,
			&position.StkCd,
			&position.Market,
			&position.Qty,
			&position.AvailableQty,
			&position.AvgPrice,
			&position.LastPrice,
			&position.HighestPrice,
			&position.InvestedAmount,
			&position.CreatedAt,
			&position.UpdatedAt,
		); err != nil {
			logger.Error("GetHoldingPositions :: error :: " + err.Error())
			return nil, err
		}
		holdingPositions = append(holdingPositions, position)
	}

	return holdingPositions, nil
}
