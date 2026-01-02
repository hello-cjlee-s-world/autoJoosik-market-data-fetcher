package repository

import (
	"autoJoosik-market-data-fetcher/internal/model"
	"autoJoosik-market-data-fetcher/pkg/logger"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UpsertTradeInfo(ctx context.Context, pool *pgxpool.Pool, entity model.TbTradeInfoLogEntity) error {
	_, err := pool.Exec(ctx, `
    INSERT INTO tb_trade_info_log (
        stk_cd, tm, cur_prc, pred_pre, pre_rt,
        pri_sel_bid_unit, pri_buy_bid_unit, cntr_trde_qty, sign,
        acc_trde_qty, acc_trde_prica, cntr_str, stex_tp, created_at
    ) VALUES (
        $1, $2, $3, $4, $5,
        $6, $7, $8, $9,
        $10, $11, $12, $13, NOW()
    )
    ON CONFLICT (stk_cd, tm, cur_prc, cntr_trde_qty) DO NOTHING
`,
		entity.StkCd,
		entity.Tm,
		entity.CurPrc,
		entity.PredPre,
		entity.PreRt,
		entity.PriSelBidUnit,
		entity.PriBuyBidUnit,
		entity.CntrTrdeQty,
		entity.Sign,
		entity.AccTrdeQty,
		entity.AccTrdePrica,
		entity.CntrStr,
		entity.StexTp,
	)

	if err != nil {
		logger.Error("UpsertTradeInfo :: error :: ", err.Error())
		return err
	}
	logger.Debug("UpsertTradeInfo :: success :: ", "stk_cd", entity.StkCd)
	return nil
}

func UpsertTradeInfoBatch(ctx context.Context, pool *pgxpool.Pool, entities []model.TbTradeInfoLogEntity) error {
	batch := &pgx.Batch{}
	for _, entity := range entities {
		batch.Queue(`
        INSERT INTO tb_trade_info_log (
            stk_cd, tm, cur_prc, pred_pre, pre_rt,
            pri_sel_bid_unit, pri_buy_bid_unit, cntr_trde_qty, sign,
            acc_trde_qty, acc_trde_prica, cntr_str, stex_tp, created_at
        ) VALUES (
            $1,$2,$3,$4,$5,
            $6,$7,$8,$9,
            $10,$11,$12,$13,NOW()
        )
        ON CONFLICT (stk_cd, tm, cur_prc, cntr_trde_qty) DO NOTHING
    `,
			entity.StkCd,
			entity.Tm,
			entity.CurPrc,
			entity.PredPre,
			entity.PreRt,
			entity.PriSelBidUnit,
			entity.PriBuyBidUnit,
			entity.CntrTrdeQty,
			entity.Sign,
			entity.AccTrdeQty,
			entity.AccTrdePrica,
			entity.CntrStr,
			entity.StexTp,
		)
	}

	br := pool.SendBatch(ctx, batch)
	defer br.Close()

	for range entities {
		if _, err := br.Exec(); err != nil {
			logger.Error("UpsertTradeInfoBatch :: error :: " + err.Error())
			return err
		}
	}
	logger.Debug("UpsertStockInfo :: success :: ")

	return br.Close()
}
