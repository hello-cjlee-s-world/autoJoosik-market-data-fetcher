package repository

import (
	"autoJoosik-market-data-fetcher/internal/model"
	"autoJoosik-market-data-fetcher/pkg/logger"
	"context"
)

func UpsertStockScore(ctx context.Context, db DB, entity model.TbStockScoreEntity) error {
	_, err := db.Exec(ctx, `
	INSERT INTO tb_stock_score(
		stk_cd, score_total, score_fundamental, score_momentum,
	    score_market, score_risk, last_price, r1, r2, r3,
	    volatility, asof_tm, meta, created_at, updated_at
	) VALUES (
	 $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,now(),now()         
	)
	ON CONFLICT (stk_cd)
	DO UPDATE SET
	  score_total = EXCLUDED.score_total,
	  score_fundamental = EXCLUDED.score_fundamental,
	  score_momentum = EXCLUDED.score_momentum,
	  score_market = EXCLUDED.score_market,
	  score_risk = EXCLUDED.score_risk,
	  last_price = EXCLUDED.last_price,
	  r1 = EXCLUDED.r1,
	  r2 = EXCLUDED.r2,
	  r3 = EXCLUDED.r3,
	  volatility = EXCLUDED.volatility,
	  asof_tm = EXCLUDED.asof_tm,
	  meta = EXCLUDED.meta,
	  created_at = now(),
	  updated_at = now();
`,
		entity.StkCd,
		entity.ScoreTotal,
		entity.ScoreFundamental,
		entity.ScoreMomentum,
		entity.ScoreMarket,
		entity.ScoreRisk,
		entity.LastPrice,
		entity.R1,
		entity.R2,
		entity.R3,
		entity.Volatility,
		entity.AsofTm,
		entity.Meta,
	)
	if err != nil {
		logger.Error("UpsertStockScore :: error :: " + err.Error())
		return err
	}
	return nil
}
