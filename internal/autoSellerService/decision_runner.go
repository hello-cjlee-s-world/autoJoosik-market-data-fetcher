package autoSellerService

import (
	"autoJoosik-market-data-fetcher/internal/model"
	"autoJoosik-market-data-fetcher/internal/repository"
	"autoJoosik-market-data-fetcher/pkg/logger"
	"context"
	"math"
	"strconv"
	"time"
)

func DecideAndExecute(ctx context.Context, pool repository.DB) error {
	decisionMutex.Lock()
	defer decisionMutex.Unlock()

	// 2보유 종목 조회
	positions, err := repository.GetHoldingPositions(ctx, pool, 0)
	if err != nil {
		return err
	}

	// ===== 매도 먼저 =====
	for _, p := range positions {
		currentPrice := p.LastPrice

		sell := ShouldSell(
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
			err := Sell(p.StkCd, p.Qty)
			if err != nil {
				return err
			}
		}
	}

	// ===== 매수 =====
	candidates, err := repository.GetBuyCandidates(ctx, pool, 0)
	if err != nil {
		return err
	}

	for _, c := range candidates {
		var bullBear = model.BullBearEntity{}
		var market = MarketState{}
		bullBear, err = repository.GetBullBearValue(ctx, pool, c.StkCd)
		market = buildMarketState(bullBear.R1, bullBear.R2, bullBear.R3, bullBear.Volatility)
		if err != nil {
			return err
		} else {

		}
		buy := ShouldBuy(
			time.Now(),
			market,
			StockState{
				StkCd:        c.StkCd,
				CurrentPrice: c.LastPrice,
				Score:        c.ScoreTotal,
			},
			c.AlreadyHolding,
			c.LastBuyTime,
			c.DailyBuyCount,
			c.CurrentHoldingCount,
			defaultBuyConstraints(),
		)
		logger.Debug("DecideAndExecute :: Finish :: " + strconv.FormatBool(buy.Do))

		if buy.Do {
			logger.Info("Buy decision", "stkCd", c.StkCd, "reason", buy.Reason)
			return Buy(c.StkCd)
		} else {
			logger.Info("Buy decision", "stkCd", c.StkCd, "reason", buy.Reason)
		}
	}

	return nil
}

func defaultBuyConstraints() BuyConstraints {
	return BuyConstraints{
		MaxHoldingCount:      10,
		MaxDailyBuyCount:     10,
		CooldownAfterBuy:     10 * time.Minute,
		AllowAddBuy:          true,
		MaxInvestPerStockPct: 0.2,
	}
}

func buildMarketState(r1, r2, r3, vol float64) MarketState {
	state := MarketState{}

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
