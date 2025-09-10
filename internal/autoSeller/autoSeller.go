package autoSeller

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type stockTick struct {
	id int
	stk_cd
	last_tic_cnt
	cur_prc
	trde_qty
	cntr_tm
	open_pric
	high_pric
	low_pric
	upd_stkpc_tp
	upd_rt
	bic_inds_tp
	sm_inds_tp
	stk_infr
	upd_stkpc_event
	pred_close_pric
}

func SellOrBy(ctx context.Context, pool *pgxpool.Pool) (string, err) {
	rows, err := pool.Query(ctx, `SELECT resources FROM possible_resources`)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var c TaskConfig
		if err := rows.Scan(&c.Name, &c.Schedule, &c.TaskType); err != nil {
			return nil, err
		}
		out = append(out, c)
	}

	return rows, err
}
