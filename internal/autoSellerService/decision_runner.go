package autoSellerService

import (
	"context"
	"time"

	"autoJoosik-market-data-fetcher/internal/repository"
	"autoJoosik-market-data-fetcher/pkg/logger"
)

func DecideAndExecute(ctx context.Context, pool repository.DB) error {
	decisionMutex.Lock()
	defer decisionMutex.Unlock()

	accountId := int64(0)

	// 1️ 시장 상태
	market, err := repository.GetMarketState(ctx, pool)
	if err != nil {
		return err
	}

	// 2️ 보유 종목 조회
	positions, err := repository.GetHoldingPositions(ctx, pool, 0)
	if err != nil {
		return err
	}

	// ===== 매도 먼저 =====
	for _, p := range positions {
		currentPrice := p.LastPrice

		sell := ShouldSell(
			market,
			Position{
				StkCd:        p.StkCd,
				Qty:          p.Qty,
				AvgPrice:     p.AvgPrice,
				HighestPrice: p.HighestPrice,
				BuyTime:      p.CreatedAt,
			},
			currentPrice,
			-1.5,
			1.0,
			-0.8,
		)

		if sell.Do {
			logger.Info("Sell decision", "stkCd", p.StkCd, "reason", sell.Reason)
			return ExecuteSell(ctx, pool, p, sell.Reason)
		}
	}

	// ===== 매수 =====
	candidates, err := repository.GetBuyCandidates(ctx, pool)
	if err != nil {
		return err
	}

	for _, c := range candidates {
		buy := ShouldBuy(
			time.Now(),
			market,
			StockState{
				StkCd:        c.StkCd,
				CurrentPrice: c.LastPrice,
				Score:        c.Score,
			},
			c.AlreadyHolding,
			c.LastBuyTime,
			c.DailyBuyCount,
			c.CurrentHoldingCount,
			DefaultBuyConstraints(),
		)

		if buy.Do {
			logger.Info("Buy decision", "stkCd", c.StkCd, "reason", buy.Reason)
			return ExecuteBuy(ctx, pool, c, buy.Reason)
		}
	}

	return nil
}
