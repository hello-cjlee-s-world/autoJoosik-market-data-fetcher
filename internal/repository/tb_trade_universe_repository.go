package repository

import (
	"autoJoosik-market-data-fetcher/internal/model"
	"autoJoosik-market-data-fetcher/pkg/logger"
	"context"
)

func GetTradeUniverse(ctx context.Context, pool DB) ([]model.TbTradeUniverseEntity, error) {
	var entities []model.TbTradeUniverseEntity

	rows, err := pool.Query(ctx, `
	SELECT 
		stk_cd,
		market,
		name,
		enabled,
		status,
		status_reason,
		created_at,
		updated_at
	FROM tb_trade_universe
	WHERE enabled = true
  		AND status = 'NORMAL'
	LIMIT 15;
`)
	if err != nil {
		logger.Error("GetTradeUniverse :: error :: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c model.TbTradeUniverseEntity
		if err := rows.Scan(
			&c.StkCd,
			&c.Market,
			&c.Name,
			&c.Enabled,
			&c.Status,
			&c.StatusReason,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			logger.Error("GetTradeUniverse :: error :: " + err.Error())
			return nil, err
		}
		entities = append(entities, c)
	}

	if err := rows.Err(); err != nil {
		logger.Error("GetTradeUniverse :: error :: " + err.Error())
		return nil, err
	}

	return entities, nil
}
