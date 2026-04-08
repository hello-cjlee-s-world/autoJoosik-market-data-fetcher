package autoSellerService

import (
	"autoJoosik-market-data-fetcher/internal/datasource"
	"autoJoosik-market-data-fetcher/internal/model"
	"autoJoosik-market-data-fetcher/internal/repository"
	"autoJoosik-market-data-fetcher/pkg/logger"
	"context"
)

func GetStockInfos() ([]model.TbStockInfoEntity, error) {
	ctx := context.Background()
	pool := datasource.GetPool()

	entity, err := repository.GetStockInfos(ctx, pool)
	if err != nil {
		logger.Error("service:GetStockInfos failed", err)
		return nil, err
	}

	return entity, nil
}
