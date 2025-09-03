package repository

import (
	"autoJoosik-market-data-fetcher/internal/model"
	"autoJoosik-market-data-fetcher/pkg/logger"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InsertStockInfo(context context.Context, pool *pgxpool.Pool, entity model.StockInfoEntity) error {
	_, err := pool.Exec(context,
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
    )`,
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

	logger.Debug("InsertStockInfo :: ", entity.StkCd, entity.StkNm)
	if err != nil {
		logger.Error("InsertStockInfo :: ", "error :: ", err.Error())
		return err
	}

	return nil
}
