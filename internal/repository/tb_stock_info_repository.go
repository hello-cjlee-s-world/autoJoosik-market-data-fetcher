package repository

import (
	"autoJoosik-market-data-fetcher/internal/model"
	"autoJoosik-market-data-fetcher/pkg/logger"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UpsertStockInfo(ctx context.Context, pool *pgxpool.Pool, entity model.TbStockInfoEntity) error {
	_, err := pool.Exec(ctx,
		`INSERT INTO tb_stock_info (
        stk_cd, stk_nm, setl_mm, fav, cap,
        flo_stk, crd_rt, oyr_hgst, oyr_lwst, mac,
        mac_wght, for_exh_rt, repl_pric, per, eps,
        roe, pbr, ev, bps, sale_amt,
        bus_pro, cup_nga, "250hgst", "250lwst", high_pric,
        open_pric, low_pric, upl_pric, lst_pric, base_pric,
        exp_cntr_pric, exp_cntr_qty, "250hgst_pric_dt", "250hgst_pric_pre_rt",
        "250lwst_pric_dt", "250lwst_pric_pre_rt", cur_prc, pre_sig, pred_pre,
        flu_rt, trde_qty, trde_pre, fav_unit, dstr_stk, dstr_rt, updated_at
    ) VALUES (
        $1, $2, $3, $4, $5,
        $6, $7, $8, $9, $10,
        $11, $12, $13, $14, $15,
        $16, $17, $18, $19, $20,
        $21, $22, $23, $24, $25,
        $26, $27, $28, $29, $30,
        $31, $32, $33, $34, $35,
        $36, $37, $38, $39, $40,
        $41, $42, $43, $44, $45, now()
    )
    ON CONFLICT (stk_cd) DO UPDATE SET
        stk_nm = EXCLUDED.stk_nm,
        setl_mm = EXCLUDED.setl_mm,
        fav = EXCLUDED.fav,
        cap = EXCLUDED.cap,
        flo_stk = EXCLUDED.flo_stk,
        crd_rt = EXCLUDED.crd_rt,
        oyr_hgst = EXCLUDED.oyr_hgst,
        oyr_lwst = EXCLUDED.oyr_lwst,
        mac = EXCLUDED.mac,
        mac_wght = EXCLUDED.mac_wght,
        for_exh_rt = EXCLUDED.for_exh_rt,
        repl_pric = EXCLUDED.repl_pric,
        per = EXCLUDED.per,
        eps = EXCLUDED.eps,
        roe = EXCLUDED.roe,
        pbr = EXCLUDED.pbr,
        ev = EXCLUDED.ev,
        bps = EXCLUDED.bps,
        sale_amt = EXCLUDED.sale_amt,
        bus_pro = EXCLUDED.bus_pro,
        cup_nga = EXCLUDED.cup_nga,
        "250hgst" = EXCLUDED."250hgst",
        "250lwst" = EXCLUDED."250lwst",
        high_pric = EXCLUDED.high_pric,
        open_pric = EXCLUDED.open_pric,
        low_pric = EXCLUDED.low_pric,
        upl_pric = EXCLUDED.upl_pric,
        lst_pric = EXCLUDED.lst_pric,
        base_pric = EXCLUDED.base_pric,
        exp_cntr_pric = EXCLUDED.exp_cntr_pric,
        exp_cntr_qty = EXCLUDED.exp_cntr_qty,
        "250hgst_pric_dt" = EXCLUDED."250hgst_pric_dt",
        "250hgst_pric_pre_rt" = EXCLUDED."250hgst_pric_pre_rt",
        "250lwst_pric_dt" = EXCLUDED."250lwst_pric_dt",
        "250lwst_pric_pre_rt" = EXCLUDED."250lwst_pric_pre_rt",
        cur_prc = EXCLUDED.cur_prc,
        pre_sig = EXCLUDED.pre_sig,
        pred_pre = EXCLUDED.pred_pre,
        flu_rt = EXCLUDED.flu_rt,
        trde_qty = EXCLUDED.trde_qty,
        trde_pre = EXCLUDED.trde_pre,
        fav_unit = EXCLUDED.fav_unit,
        dstr_stk = EXCLUDED.dstr_stk,
        dstr_rt = EXCLUDED.dstr_rt,
        updated_at = now()
    `,
		// values 그대로
		entity.StkCd,
		entity.StkNm,
		entity.SetlMm,
		entity.Fav,
		entity.Cap,
		entity.FloStk,
		entity.CrdRt,
		entity.OyrHgst,
		entity.OyrLwst,
		entity.Mac,
		entity.MacWght,
		entity.ForExhRt,
		entity.ReplPric,
		entity.Per,
		entity.Eps,
		entity.Roe,
		entity.Pbr,
		entity.Ev,
		entity.Bps,
		entity.SaleAmt,
		entity.BusPro,
		entity.CupNga,
		entity.Hgst250,
		entity.Lwst250,
		entity.HighPric,
		entity.OpenPric,
		entity.LowPric,
		entity.UplPric,
		entity.LstPric,
		entity.BasePric,
		entity.ExpCntrPric,
		entity.ExpCntrQty,
		entity.Hgst250PricDt,
		entity.Hgst250PreRt,
		entity.Lwst250PricDt,
		entity.Lwst250PreRt,
		entity.CurPrc,
		entity.PreSig,
		entity.PredPre,
		entity.FluRt,
		entity.TrdeQty,
		entity.TrdePre,
		entity.FavUnit,
		entity.DstrStk,
		entity.DstrRt,
	)

	if err != nil {
		logger.Error("UpsertStockInfo :: error :: ", "stk_cd", entity.StkCd)
		return err
	}
	logger.Debug("UpsertStockInfo :: success :: ", "stk_cd", entity.StkCd)
	return nil
}

func UpsertStockInfoBatch(ctx context.Context, pool *pgxpool.Pool, entities []model.TbStockInfoEntity) error {
	if len(entities) == 0 {
		return nil
	}

	batch := &pgx.Batch{}

	for _, entity := range entities {
		batch.Queue(`
			INSERT INTO tb_stock_info (
				stk_cd, stk_nm, setl_mm, fav, cap,
				flo_stk, crd_rt, oyr_hgst, oyr_lwst, mac,
				mac_wght, for_exh_rt, repl_pric, per, eps,
				roe, pbr, ev, bps, sale_amt,
				bus_pro, cup_nga, "250hgst", "250lwst", high_pric,
				open_pric, low_pric, upl_pric, lst_pric, base_pric,
				exp_cntr_pric, exp_cntr_qty, "250hgst_pric_dt", "250hgst_pric_pre_rt",
				"250lwst_pric_dt", "250lwst_pric_pre_rt", cur_prc, pre_sig, pred_pre,
				flu_rt, trde_qty, trde_pre, fav_unit, dstr_stk, dstr_rt, updated_at
			) VALUES (
				$1,$2,$3,$4,$5,
				$6,$7,$8,$9,$10,
				$11,$12,$13,$14,$15,
				$16,$17,$18,$19,$20,
				$21,$22,$23,$24,$25,
				$26,$27,$28,$29,$30,
				$31,$32,$33,$34,$35,
				$36,$37,$38,$39,$40,
				$41,$42,$43,$44,$45, now()
			)
			ON CONFLICT (stk_cd) DO UPDATE SET
				stk_nm = EXCLUDED.stk_nm,
				setl_mm = EXCLUDED.setl_mm,
				fav = EXCLUDED.fav,
				cap = EXCLUDED.cap,
				flo_stk = EXCLUDED.flo_stk,
				crd_rt = EXCLUDED.crd_rt,
				oyr_hgst = EXCLUDED.oyr_hgst,
				oyr_lwst = EXCLUDED.oyr_lwst,
				mac = EXCLUDED.mac,
				mac_wght = EXCLUDED.mac_wght,
				for_exh_rt = EXCLUDED.for_exh_rt,
				repl_pric = EXCLUDED.repl_pric,
				per = EXCLUDED.per,
				eps = EXCLUDED.eps,
				roe = EXCLUDED.roe,
				pbr = EXCLUDED.pbr,
				ev = EXCLUDED.ev,
				bps = EXCLUDED.bps,
				sale_amt = EXCLUDED.sale_amt,
				bus_pro = EXCLUDED.bus_pro,
				cup_nga = EXCLUDED.cup_nga,
				"250hgst" = EXCLUDED."250hgst",
				"250lwst" = EXCLUDED."250lwst",
				high_pric = EXCLUDED.high_pric,
				open_pric = EXCLUDED.open_pric,
				low_pric = EXCLUDED.low_pric,
				upl_pric = EXCLUDED.upl_pric,
				lst_pric = EXCLUDED.lst_pric,
				base_pric = EXCLUDED.base_pric,
				exp_cntr_pric = EXCLUDED.exp_cntr_pric,
				exp_cntr_qty = EXCLUDED.exp_cntr_qty,
				"250hgst_pric_dt" = EXCLUDED."250hgst_pric_dt",
				"250hgst_pric_pre_rt" = EXCLUDED."250hgst_pric_pre_rt",
				"250lwst_pric_dt" = EXCLUDED."250lwst_pric_dt",
				"250lwst_pric_pre_rt" = EXCLUDED."250lwst_pric_pre_rt",
				cur_prc = EXCLUDED.cur_prc,
				pre_sig = EXCLUDED.pre_sig,
				pred_pre = EXCLUDED.pred_pre,
				flu_rt = EXCLUDED.flu_rt,
				trde_qty = EXCLUDED.trde_qty,
				trde_pre = EXCLUDED.trde_pre,
				fav_unit = EXCLUDED.fav_unit,
				dstr_stk = EXCLUDED.dstr_stk,
				dstr_rt = EXCLUDED.dstr_rt,
				updated_at = now()
		`,
			entity.StkCd,
			entity.StkNm,
			entity.SetlMm,
			entity.Fav,
			entity.Cap,
			entity.FloStk,
			entity.CrdRt,
			entity.OyrHgst,
			entity.OyrLwst,
			entity.Mac,
			entity.MacWght,
			entity.ForExhRt,
			entity.ReplPric,
			entity.Per,
			entity.Eps,
			entity.Roe,
			entity.Pbr,
			entity.Ev,
			entity.Bps,
			entity.SaleAmt,
			entity.BusPro,
			entity.CupNga,
			entity.Hgst250,
			entity.Lwst250,
			entity.HighPric,
			entity.OpenPric,
			entity.LowPric,
			entity.UplPric,
			entity.LstPric,
			entity.BasePric,
			entity.ExpCntrPric,
			entity.ExpCntrQty,
			entity.Hgst250PricDt,
			entity.Hgst250PreRt,
			entity.Lwst250PricDt,
			entity.Lwst250PreRt,
			entity.CurPrc,
			entity.PreSig,
			entity.PredPre,
			entity.FluRt,
			entity.TrdeQty,
			entity.TrdePre,
			entity.FavUnit,
			entity.DstrStk,
			entity.DstrRt,
		)
	}

	br := pool.SendBatch(ctx, batch)
	defer br.Close()

	for range entities {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}

	return nil
}

// 기업 점수 계산 위한 컬럼만 select
func GetStockFundamental(ctx context.Context, pool DB, stkCd string) (model.TbStockInfoEntity, error) {
	var entity model.TbStockInfoEntity
	err := pool.QueryRow(ctx, `
SELECT
  per,
  roe,
  pbr,
  eps,
  for_exh_rt,
  cap
FROM tb_stock_info
WHERE stk_cd = $1;
`, stkCd).Scan(&entity.Per, &entity.Roe, &entity.Pbr, &entity.Eps, &entity.ForExhRt, &entity.Cap)
	if err != nil {
		logger.Error("GetStockFundamental :: error :: " + err.Error())
	}

	return entity, nil
}
