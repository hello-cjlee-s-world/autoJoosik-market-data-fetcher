package autoSellerService

import (
	"autoJoosik-market-data-fetcher/internal/repository"
	"context"
)

func DecideAndExecute(ctx context.Context, pool repository.DB) error {
	//decisionMutex.Lock()
	//defer decisionMutex.Unlock()
	//
	//// 2보유 종목 조회
	//positions, err := repository.GetHoldingPositions(ctx, pool, 0)
	//if err != nil {
	//	return err
	//}
	//
	//// ===== 매도 먼저 =====
	//for _, p := range positions {
	//	currentPrice := p.LastPrice
	//
	//	sell := ShouldSell(
	//		Position{
	//			StkCd:        p.StkCd,
	//			Qty:          p.Qty,
	//			AvgPrice:     p.AvgPrice,
	//			HighestPrice: p.HighestPrice,
	//			BuyTime:      p.CreatedAt,
	//		},
	//		currentPrice,
	//		-1.5,
	//		1.0,
	//		-0.8,
	//	)
	//
	//	if sell.Do {
	//		logger.Info("Sell decision", "stkCd", p.StkCd, "reason", sell.Reason)
	//		Sell(p.StkCd, p.Qty)
	//	}
	//}
	//
	//// ===== 매수 =====
	////candidates, err := repository.GetBuyCandidates(ctx, pool)
	//var candidates []Candidate
	//if err != nil {
	//	return err
	//}
	//
	//for _, c := range candidates {
	//	var bullBear = model.BullBearEntity{}
	//	var market = MarketState{}
	//	bullBear, err = repository.GetBullBearValue(ctx, pool, c.StkCd)
	//	market = repository.BuildMarketState(bullBear.R1, bullBear.R2, bullBear.R3, bullBear.Volatility)
	//	if err != nil {
	//		return err
	//	} else {
	//
	//	}
	//
	//	buy := ShouldBuy(
	//		time.Now(),
	//		market,
	//		StockState{
	//			StkCd:        c.StkCd,
	//			CurrentPrice: c.LastPrice,
	//			Score:        c.Score,
	//		},
	//		c.AlreadyHolding,
	//		c.LastBuyTime,
	//		c.DailyBuyCount,
	//		c.CurrentHoldingCount,
	//		DefaultBuyConstraints(),
	//	)
	//
	//	if buy.Do {
	//		logger.Info("Buy decision", "stkCd", c.StkCd, "reason", buy.Reason)
	//		return ExecuteBuy(ctx, pool, c, buy.Reason)
	//	}
	//}

	return nil
}
