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
		  invested_amount,
		  today_buy_qty, today_sell_qty,
		  status, last_eval_at,
		  created_at, updated_at
		) VALUES (
		  $1, $2, $3, $4, $7::text,
		  $5, $5, $6, $6, $6,
		  ($5::numeric * $6::numeric),
		  CASE WHEN $7::text = 'B' THEN $5 ELSE 0 END,
		  CASE WHEN $7::text = 'S' THEN $5 ELSE 0 END,
		  $8, NOW(),
		  NOW(), NOW()
		)
		ON CONFLICT (account_id, stk_cd, market)
		DO UPDATE SET
		  user_id = EXCLUDED.user_id,

		  qty = CASE
			WHEN $7::text = 'B' THEN tb_virtual_asset.qty + EXCLUDED.qty
			WHEN $7::text = 'S' THEN tb_virtual_asset.qty - EXCLUDED.qty
			ELSE tb_virtual_asset.qty
		  END,

		  available_qty = CASE
			WHEN $7::text = 'B' THEN tb_virtual_asset.available_qty + EXCLUDED.qty
			WHEN $7::text = 'S' THEN tb_virtual_asset.available_qty - EXCLUDED.qty
			ELSE tb_virtual_asset.available_qty
		  END,

		  avg_price = CASE
			WHEN $7::text = 'B'
			 AND (tb_virtual_asset.qty + EXCLUDED.qty) > 0
			  THEN (tb_virtual_asset.invested_amount + EXCLUDED.invested_amount)
				   / (tb_virtual_asset.qty + EXCLUDED.qty)
			ELSE tb_virtual_asset.avg_price
		  END,

		  last_price = EXCLUDED.last_price,
		     
		  highest_price = CASE
		     WHEN $7::text = 'B'
			  AND tb_virtual_asset.highest_price < EXCLUDED.avg_price
			   THEN EXCLUDED.avg_price
		     WHEN $7::text = 'S'
		       THEN tb_virtual_asset.highest_price
		    END,
			 

		  invested_amount = CASE
			WHEN $7::text = 'B' THEN tb_virtual_asset.invested_amount + EXCLUDED.invested_amount
			WHEN $7::text = 'S' THEN tb_virtual_asset.invested_amount - EXCLUDED.invested_amount
			ELSE tb_virtual_asset.invested_amount
		  END,

		  today_buy_qty  = tb_virtual_asset.today_buy_qty  + EXCLUDED.today_buy_qty,
		  today_sell_qty = tb_virtual_asset.today_sell_qty + EXCLUDED.today_sell_qty,

		  status       = EXCLUDED.status,
		  last_eval_at = NOW(),
		  updated_at   = NOW();
		`,
		entity.UserId,
		entity.AccountId,
		entity.StkCd,
		entity.Market,
		entity.Qty,
		entity.AvgPrice,
		entity.PositionSide, // 여기서는 "side"로만 사용 (B/S)
		entity.Status,
	)
	if err != nil {
		logger.Error("UpsertVirtualAsset :: error :: ", err)
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

func GetHoldingPositions(ctx context.Context, db DB, accountId int64) ([]model.HoldingPosition, error) {
	var virtualAssetEntities []model.TbVirtualAssetEntity
	var holdingPositions []model.HoldingPosition

	err := db.QueryRow(ctx, `
	SELECT * 
	FROM tb_virtual_asset
	WHERE account_id = $1
`, accountId).Scan(&virtualAssetEntities)

	if err != nil {
		logger.Error("GetHoldingPositions :: error :: " + err.Error())
	}

	if virtualAssetEntities != nil && len(virtualAssetEntities) > 0 {
		for _, virtualAssetEntity := range virtualAssetEntities {
			var position = model.HoldingPosition{
				AccountId:    virtualAssetEntity.AccountId,
				UserId:       virtualAssetEntity.UserId,
				StkCd:        virtualAssetEntity.StkCd,
				Market:       virtualAssetEntity.Market,
				Qty:          virtualAssetEntity.Qty,
				AvailableQty: virtualAssetEntity.AvailableQty,
				AvgPrice:     virtualAssetEntity.AvgPrice,
				LastPrice:    virtualAssetEntity.LastPrice,
				HighestPrice: virtualAssetEntity.HighestPrice,
				InvestedAmt:  virtualAssetEntity.InvestedAmount,
				CreatedAt:    virtualAssetEntity.CreatedAt,
				UpdatedAt:    virtualAssetEntity.UpdatedAt,
			}
			holdingPositions = append(holdingPositions, position)
		}
	}

	return holdingPositions, nil
}
