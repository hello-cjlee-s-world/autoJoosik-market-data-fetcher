package repository

import (
	"autoJoosik-market-data-fetcher/internal/model"
	"autoJoosik-market-data-fetcher/pkg/logger"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
)

func UpsertStockInfo(ctx context.Context, pool *pgxpool.Pool, entity model.StockInfoEntity) error {
	_, err := pool.Exec(ctx,
		`INSERT INTO stock_info (
        stk_cd, stk_nm, setl_mm, fav, cap,
        flo_stk, crd_rt, oyr_hgst, oyr_lwst, mac,
        mac_wght, for_exh_rt, repl_pric, per, eps,
        roe, pbr, ev, bps, sale_amt,
        bus_pro, cup_nga, "250hgst", "250lwst", high_pric,
        open_pric, low_pric, upl_pric, lst_pric, base_pric,
        exp_cntr_pric, exp_cntr_qty, "250hgst_pric_dt", "250hgst_pric_pre_rt",
        "250lwst_pric_dt", "250lwst_pric_pre_rt", cur_prc, pre_sig, pred_pre,
        flu_rt, trde_qty, trde_pre, fav_unit, dstr_stk, dstr_rt
    ) VALUES (
        $1, $2, $3, $4, $5,
        $6, $7, $8, $9, $10,
        $11, $12, $13, $14, $15,
        $16, $17, $18, $19, $20,
        $21, $22, $23, $24, $25,
        $26, $27, $28, $29, $30,
        $31, $32, $33, $34, $35,
        $36, $37, $38, $39, $40,
        $41, $42, $43, $44, $45
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
        dstr_rt = EXCLUDED.dstr_rt
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

func UpsertStockInfoBatch(ctx context.Context, pool *pgxpool.Pool, entities []model.StockInfoEntity) error {
	if len(entities) == 0 {
		return nil
	}

	// VALUES ($1, $2, ...), ($n, $n+1, ...)
	valueStrings := make([]string, 0, len(entities))
	valueArgs := make([]interface{}, 0, len(entities)*45) // 컬럼 개수만큼 곱해야 함

	for i, e := range entities {
		base := i * 45
		valueStrings = append(valueStrings,
			fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,"+
				"$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,"+
				"$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,"+
				"$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,"+
				"$%d,$%d,$%d,$%d,$%d)",
				base+1, base+2, base+3, base+4, base+5,
				base+6, base+7, base+8, base+9, base+10,
				base+11, base+12, base+13, base+14, base+15,
				base+16, base+17, base+18, base+19, base+20,
				base+21, base+22, base+23, base+24, base+25,
				base+26, base+27, base+28, base+29, base+30,
				base+31, base+32, base+33, base+34, base+35,
				base+36, base+37, base+38, base+39, base+40,
				base+41, base+42, base+43, base+44, base+45))
		valueArgs = append(valueArgs,
			e.StkCd, e.StkNm, e.SetlMm, e.Fav, e.Cap,
			e.FloStk, e.CrdRt, e.OyrHgst, e.OyrLwst, e.Mac,
			e.MacWght, e.ForExhRt, e.ReplPric, e.Per, e.Eps,
			e.Roe, e.Pbr, e.Ev, e.Bps, e.SaleAmt,
			e.BusPro, e.CupNga, e.Hgst250, e.Lwst250, e.HighPric,
			e.OpenPric, e.LowPric, e.UplPric, e.LstPric, e.BasePric,
			e.ExpCntrPric, e.ExpCntrQty, e.Hgst250PricDt, e.Hgst250PreRt,
			e.Lwst250PricDt, e.Lwst250PreRt, e.CurPrc, e.PreSig, e.PredPre,
			e.FluRt, e.TrdeQty, e.TrdePre, e.FavUnit, e.DstrStk, e.DstrRt,
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO stock_info (
			stk_cd, stk_nm, setl_mm, fav, cap,
			flo_stk, crd_rt, oyr_hgst, oyr_lwst, mac,
			mac_wght, for_exh_rt, repl_pric, per, eps,
			roe, pbr, ev, bps, sale_amt,
			bus_pro, cup_nga, "250hgst", "250lwst", high_pric,
			open_pric, low_pric, upl_pric, lst_pric, base_pric,
			exp_cntr_pric, exp_cntr_qty, "250hgst_pric_dt", "250hgst_pric_pre_rt",
			"250lwst_pric_dt", "250lwst_pric_pre_rt", cur_prc, pre_sig, pred_pre,
			flu_rt, trde_qty, trde_pre, fav_unit, dstr_stk, dstr_rt
		) VALUES %s
		ON CONFLICT (stk_cd) DO UPDATE SET
			stk_nm=EXCLUDED.stk_nm, setl_mm=EXCLUDED.setl_mm, fav=EXCLUDED.fav, cap=EXCLUDED.cap,
			flo_stk=EXCLUDED.flo_stk, crd_rt=EXCLUDED.crd_rt, oyr_hgst=EXCLUDED.oyr_hgst, oyr_lwst=EXCLUDED.oyr_lwst, mac=EXCLUDED.mac,
			mac_wght=EXCLUDED.mac_wght, for_exh_rt=EXCLUDED.for_exh_rt, repl_pric=EXCLUDED.repl_pric, per=EXCLUDED.per, eps=EXCLUDED.eps,
			roe=EXCLUDED.roe, pbr=EXCLUDED.pbr, ev=EXCLUDED.ev, bps=EXCLUDED.bps, sale_amt=EXCLUDED.sale_amt,
			bus_pro=EXCLUDED.bus_pro, cup_nga=EXCLUDED.cup_nga, "250hgst"=EXCLUDED."250hgst", "250lwst"=EXCLUDED."250lwst", high_pric=EXCLUDED.high_pric,
			open_pric=EXCLUDED.open_pric, low_pric=EXCLUDED.low_pric, upl_pric=EXCLUDED.upl_pric, lst_pric=EXCLUDED.lst_pric, base_pric=EXCLUDED.base_pric,
			exp_cntr_pric=EXCLUDED.exp_cntr_pric, exp_cntr_qty=EXCLUDED.exp_cntr_qty, "250hgst_pric_dt"=EXCLUDED."250hgst_pric_dt", "250hgst_pric_pre_rt"=EXCLUDED."250hgst_pric_pre_rt",
			"250lwst_pric_dt"=EXCLUDED."250lwst_pric_dt", "250lwst_pric_pre_rt"=EXCLUDED."250lwst_pric_pre_rt", cur_prc=EXCLUDED.cur_prc, pre_sig=EXCLUDED.pre_sig, pred_pre=EXCLUDED.pred_pre,
			flu_rt=EXCLUDED.flu_rt, trde_qty=EXCLUDED.trde_qty, trde_pre=EXCLUDED.trde_pre, fav_unit=EXCLUDED.fav_unit, dstr_stk=EXCLUDED.dstr_stk, dstr_rt=EXCLUDED.dstr_rt
	`, strings.Join(valueStrings, ","))

	_, err := pool.Exec(ctx, query, valueArgs...)
	return err
}
