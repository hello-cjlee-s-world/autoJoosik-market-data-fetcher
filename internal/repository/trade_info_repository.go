package repository

import (
	"autoJoosik-market-data-fetcher/internal/autoSellerService"
	"autoJoosik-market-data-fetcher/internal/model"
	"autoJoosik-market-data-fetcher/pkg/logger"
	"context"
	"math"
)

func GetBullBearValue(ctx context.Context, db DB, stkCd string) (model.BullBearEntity, error) {
	entity := model.BullBearEntity{}

	err := db.QueryRow(ctx, `SELECT
	  (cur_prc - LAG(cur_prc, 1) OVER w) / LAG(cur_prc, 1) OVER w * 100 AS r1,
	  (cur_prc - LAG(cur_prc, 10) OVER w) / LAG(cur_prc, 10) OVER w * 100 AS r2,
	  (cur_prc - LAG(cur_prc, 30) OVER w) / LAG(cur_prc, 30) OVER w * 100 AS r3
	  
	FROM (
	  SELECT time, cur_prc
	  FROM trade_info_log
	  WHERE stk_cd = $1
	  ORDER BY time DESC
	  LIMIT 25
	) t
	WINDOW w AS (ORDER BY time DESC)
	LIMIT 1;
`, stkCd).Scan(&entity)

	if err != nil {
		logger.Error("getBullBearValue :: error :: " + err.Error())
		return entity, err
	}

	// 변동성도 같이 계산
	var volatility float64

	err = db.QueryRow(ctx, `
	SELECT volatility
	FROM view_volatility
	WHERE stk_cd = $1
`, stkCd).Scan(&volatility)

	if err != nil {
		logger.Error("getBullBearValue :: error :: " + err.Error())
		return entity, err
	} else {
		entity.Volatility = volatility
	}

	return entity, nil
}

func BuildMarketState(r1, r2, r3, vol float64) autoSellerService.MarketState {
	state := autoSellerService.MarketState{}

	state.IsBull = r1 > 0 && r2 > 0 && r3 > 0
	state.IsBear = r1 < 0 && r2 < 0

	state.Volatility = vol
	state.IndexChange = r1

	if math.Abs(r1) >= 2.0 || vol >= 2.5 {
		state.IsEmergency = true
		state.Reason = "market_shock"
	} else {
		state.IsEmergency = false
		state.Reason = "normal"
	}

	return state
}
