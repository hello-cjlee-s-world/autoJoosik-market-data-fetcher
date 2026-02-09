package autoSellerService

import (
	"autoJoosik-market-data-fetcher/internal/datasource"
	"autoJoosik-market-data-fetcher/internal/kiwoomApi"
	"autoJoosik-market-data-fetcher/internal/model"
	"autoJoosik-market-data-fetcher/internal/repository"
	"autoJoosik-market-data-fetcher/internal/utils"
	"autoJoosik-market-data-fetcher/pkg/logger"
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func Sell(stkCd string, qty float64) error {
	rst, err := kiwoomApi.GetOrderBookLog(stkCd)
	if !utils.IsTradableTime(time.Now()) {
		return fmt.Errorf("not available trade time")
	}
	if err == nil {
		//  주식 거래 예시 트랜잭션으로 묶기
		ctx := context.Background()
		pool := datasource.GetPool()
		tx, _ := pool.Begin(ctx)

		orderBookEntity := model.ToOrderBookLogEntity(rst)

		// 계좌 내 거래 가능 수량
		availableQty, err := repository.GetAvailableAssetsByAccountAndStkCd(ctx, tx, 0, stkCd)
		if err != nil {
			logger.Error("Sell :: error :: stkCd" + stkCd + " " + err.Error())
			return err
		}
		if availableQty <= 0 {
			logger.Info("Sell :: 거래 가능한 수량이 없습니다.")
		} else if availableQty < qty && availableQty > 0 {
			qty = availableQty
			logger.Info("Sell :: 매도 가능 수량 부족, 가능한 수량만 매도합니다.")
		}

		price, _ := strconv.ParseFloat(orderBookEntity.BuyFprBid, 64)

		// !!insert 주문 정보(상태)
		remainingQty := 0.0 // 내 주문 기준 남은 수량
		userId := int64(0)
		accountID := int64(0)
		clientOrderID := fmt.Sprintf("%d_%d_%d_%s",
			accountID,
			time.Now().UnixNano(),
			rand.Intn(1000),
			"sell",
		) // 유니크한 주문 아이디
		virtualOrderEntity := model.TbVirtualOrder{
			UserID:       userId,
			AccountID:    accountID,
			StkCd:        stkCd,
			Market:       "KOSPI",
			Side:         "A",
			OrderType:    "MARKET",
			TimeInForce:  "DAY",
			Price:        price,
			Qty:          qty,
			FilledQty:    qty,
			RemainingQty: remainingQty,
			Status:       "FILLED", //'NEW','OPEN','PARTIAL','FILLED','CANCELED','REJECTED'

			ClientOrderID: clientOrderID, // 클라이언트에서 부여하는 주문 ID (선택)
			Reason:        "",            //-- 거절/취소 사유 등 (옵션)
		}
		orderId, err := repository.InsertOrder(ctx, tx, virtualOrderEntity)
		if err != nil {
			fmt.Println("InsertOrderLog", err.Error())
			return err
		}

		// !!insert 거래 로그 (원래는 주문 체결 후 insert 지만 가상이라 바로 체결함에 따라 바로 insert)
		virtualTradeLogEntity := model.TbVirtualTradeLog{
			OrderID:      orderId,
			UserID:       0,
			AccountID:    0,
			StkCd:        stkCd,
			Market:       "KOSPI",
			Side:         "S",
			FilledQty:    qty,
			FilledPrice:  price,
			FilledAmount: qty * price,
			FeeAmount:    0.1,
			TaxAmount:    0.1,
		}

		_, err = repository.InsertTradeLog(ctx, tx, virtualTradeLogEntity)
		if err != nil {
			fmt.Println("InsertTradeLog", err.Error())
			return err
		}

		// !!upsert 가상 자산 테이블 에
		virtualAssetEntity := model.TbVirtualAssetEntity{
			UserId:       0,
			AccountId:    0,
			StkCd:        stkCd,
			Market:       "KOSPI",
			PositionSide: "S",
			Qty:          qty,
			AvgPrice:     price,
			Status:       "ACTIVE",
		}
		err = repository.UpsertVirtualAsset(ctx, tx, virtualAssetEntity)
		if err != nil {
			fmt.Println("UpsertVirtualAsset", err.Error())
			return err
		}

		// !!update 가상 계좌 테이블 거래 가능 금액
		virtualAccountEntity := model.TbVirtualAccountEntity{
			AccountId: 0,
		}
		err = repository.UpdateVirtualAccount(ctx, tx, virtualAccountEntity, "SELL", int64(price*qty))
		// 트랜잭션으로 묶어서 commit
		if err := tx.Commit(ctx); err != nil {
			logger.Error("Sell :: error ::" + err.Error())
			return err
		} else {
			logger.Info("Sell :: success :: accountId=" + fmt.Sprintf(strconv.FormatInt(accountID, 10)) + ", stkCd=" + stkCd)
		}
	}
	return err
}
