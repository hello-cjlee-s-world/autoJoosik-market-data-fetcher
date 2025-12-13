package repository

import (
	"autoJoosik-market-data-fetcher/internal/model"
	"autoJoosik-market-data-fetcher/pkg/logger"
	"context"
)

func InsertVirtualAsset(ctx context.Context, db DB, entity model.TbVirtualAssetEntity) error {
	_, err := db.Exec(ctx, `
		INSERT INTO tb_virtual_asset (
		  user_id, account_id, stk_cd, market, position_side,
		  qty, available_qty, avg_price, last_price,
		  invested_amount,
		  today_buy_qty, today_sell_qty,
		  status, last_eval_at,
		  created_at, updated_at
		) VALUES (
		  $1, $2, $3, $4, $5,
		  $6, $6, $7, $7,
		  ($6::numeric * $7::numeric),
		  $6, 0,
		  $8, NOW(),
		  NOW(), NOW()
		)
		ON CONFLICT (account_id, stk_cd, market, position_side)
		DO UPDATE SET
		  user_id = EXCLUDED.user_id,
		
		  -- 수량 증가
		  qty = tb_virtual_asset.qty + EXCLUDED.qty,
		
		  -- 매수 직후엔 보통 매도가능수량도 같이 증가(락/미체결 없다는 가정)
		  available_qty = tb_virtual_asset.available_qty + EXCLUDED.qty,
		
		  -- 평단 재계산: (기존투자금 + 이번매수금) / (기존수량 + 이번수량)
		  avg_price = CASE
			WHEN (tb_virtual_asset.qty + EXCLUDED.qty) = 0 THEN 0
			ELSE (tb_virtual_asset.invested_amount + EXCLUDED.invested_amount)
				 / (tb_virtual_asset.qty + EXCLUDED.qty)
		  END,
		
		  -- 마지막 체결가/마지막 가격
		  last_price = EXCLUDED.last_price,
		
		  -- 투자금 누적
		  invested_amount = tb_virtual_asset.invested_amount + EXCLUDED.invested_amount,
		
		  -- 오늘 매수량 누적
		  today_buy_qty = tb_virtual_asset.today_buy_qty + EXCLUDED.today_buy_qty,
		
		  status = EXCLUDED.status,
		  last_eval_at = NOW(),
		  updated_at = NOW();
		`,
		entity.UserId,
		entity.AccountId,
		entity.StkCd,
		entity.Market,
		entity.PositionSide,
		entity.Qty,
		entity.AvgPrice,
		entity.Status,
	)
	if err != nil {
		logger.Error("InsertVirtualAsset :: error :: ", err)
		return err
	}
	logger.Debug("InsertTradeLog :: success :: ", "stk_cd", entity.StkCd)
	return nil
}
