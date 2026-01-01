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

func GetBuyCandidates(ctx context.Context, db DB, accountId int64) ([]model.CandidateEntity, error) {
	rows, err := db.Query(ctx, `
WITH candidates AS (
  SELECT stk_cd, last_price, score_total
  FROM tb_stock_score
  ORDER BY score_total DESC
  LIMIT 20
),
holding AS (
  SELECT stk_cd, (qty > 0) AS already_holding
  FROM tb_virtual_asset
  WHERE account_id = $1
),
today_buy AS (
  SELECT stk_cd, COUNT(*) AS daily_buy_count, MAX(created_at) AS last_buy_time
  FROM tb_virtual_trade_log
  WHERE account_id = $1
    AND side = 'B'
    AND created_at >= date_trunc('day', now())
  GROUP BY stk_cd
),
holding_cnt AS (
  SELECT COUNT(*) AS cnt
  FROM tb_virtual_asset
  WHERE account_id = $1 AND qty > 0
)
SELECT
  c.stk_cd,
  c.last_price,
  c.score_total,
  COALESCE(h.already_holding, false) AS already_holding,
  COALESCE(t.last_buy_time, '1970-01-01'::timestamptz) AS last_buy_time,
  COALESCE(t.daily_buy_count, 0) AS daily_buy_count,
  (SELECT cnt FROM holding_cnt) AS current_holding_count
FROM candidates c
LEFT JOIN holding h ON h.stk_cd = c.stk_cd
LEFT JOIN today_buy t ON t.stk_cd = c.stk_cd
ORDER BY c.score_total DESC;
`, accountId)
	if err != nil {
		logger.Error("GetBuyCandidates :: error :: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	var out []model.CandidateEntity
	for rows.Next() {
		var c model.CandidateEntity
		if err := rows.Scan(
			&c.StkCd,
			&c.LastPrice,
			&c.ScoreTotal,
			&c.AlreadyHolding,
			&c.LastBuyTime,
			&c.DailyBuyCount,
			&c.CurrentHoldingCount,
		); err != nil {
			logger.Error("GetBuyCandidates :: error :: " + err.Error())
			return nil, err
		}
		out = append(out, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
