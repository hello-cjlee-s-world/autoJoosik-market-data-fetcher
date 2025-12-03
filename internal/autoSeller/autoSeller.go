package autoSeller

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type stockTick struct {
	Id              int
	stkCd           string
	lastTicCnt      int
	curPrc          int
	trdeQty         int
	cntrTm          string
	openPric        int
	highPric        int
	lowPric         int
	updStkpcTp      int
	updRt           int
	bicIndsTp       string
	smIndsTp        string
	stk_infr        string
	upd_stkpc_event string
	pred_close_pric int
}

func SellOrBy(ctx context.Context, pool *pgxpool.Pool) (string, err) {
	rows, err := pool.Query(ctx, `SELECT resources FROM possible_resources`)

	if err != nil {
		return nil, err
	}
	out := []stockTick{}

	for rows.Next() {
		t := stockTick{}
		if err := rows.Scan(&t.Name, &c.Schedule, &tickList); err != nil {
			return nil, err
		}
		out = append(out, c)
	}

	return rows, err
}
