package repository

import (
	"autoJoosik-market-data-fetcher/internal/model"
	"autoJoosik-market-data-fetcher/pkg/logger"
	"context"
)

func GetBullBearValue(ctx context.Context, db DB, stkCd string) (model.BullBearEntity, error) {
	var entity model.BullBearEntity

	err := db.QueryRow(ctx, `
  SELECT
  COALESCE((cur_prc - LAG(cur_prc, 1)  OVER w) / NULLIF(LAG(cur_prc, 1)  OVER w, 0) * 100, 0) AS r1,
  COALESCE((cur_prc - LAG(cur_prc, 10) OVER w) / NULLIF(LAG(cur_prc, 10) OVER w, 0) * 100, 0) AS r2,
  COALESCE((cur_prc - LAG(cur_prc, 30) OVER w) / NULLIF(LAG(cur_prc, 30) OVER w, 0) * 100, 0) AS r3
  FROM (
    SELECT tm, cur_prc
    FROM tb_trade_info_log
    WHERE stk_cd = $1
    ORDER BY tm DESC
    LIMIT 40
  ) t
  WINDOW w AS (ORDER BY tm DESC)
  LIMIT 1;
`, stkCd).Scan(&entity.R1, &entity.R2, &entity.R3)

	if err != nil {
		logger.Error("getBullBearValue :: error 1 :: " + err.Error())
		return entity, err
	}

	// 변동성도 같이 계산
	err = db.QueryRow(ctx, `
		SELECT fn_volatility($1, 300)
	`, stkCd).Scan(&entity.Volatility)

	if err != nil {
		logger.Error("getBullBearValue :: error 2 :: " + err.Error())
		return entity, err
	}

	return entity, nil
}
