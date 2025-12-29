package autoSellerService

import (
	"autoJoosik-market-data-fetcher/internal/datasource"
	"autoJoosik-market-data-fetcher/internal/kiwoomApi"
	"autoJoosik-market-data-fetcher/internal/model"
	"autoJoosik-market-data-fetcher/internal/repository"
	"autoJoosik-market-data-fetcher/pkg/logger"
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func Buy() {
	stkCd := "005930"
	rst, err := kiwoomApi.GetOrderBookLog(stkCd)
	if err == nil {
		//  주식 거래 예시 트랜잭션으로 묶기
		ctx := context.Background()
		pool := datasource.GetPool()
		tx, err := pool.Begin(ctx)

		orderBookEntity := model.ToOrderBookLogEntity(rst)

		price, _ := strconv.ParseFloat(orderBookEntity.SelFprBid, 64)
		//remainingQty, _ := strconv.ParseFloat(orderBookEntity.SelFprReq, 64)

		// !!insert 주문 정보(상태)
		qty := 1.0          // 내가 사고 싶은 수량
		remainingQty := 0.0 // 내 주문 기준 남은 수량
		userId := int64(0)
		accountID := int64(0)
		clientOrderID := fmt.Sprintf("%d_%d_%d",
			accountID,
			time.Now().UnixNano(),
			rand.Intn(1000),
		) // 유니크한 주문 아이디
		virtualOrderEntity := model.TbVirtualOrder{
			UserID:       userId,
			AccountID:    accountID,
			StkCd:        stkCd,
			Market:       "KOSPI",
			Side:         "B",
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
		}

		// !!insert 거래 로그 (원래는 주문 체결 후 insert 지만 가상이라 바로 체결함에 따라 바로 insert)
		virtualTradeLogEntity := model.TbVirtualTradeLog{
			OrderID:      orderId,
			UserID:       0,
			AccountID:    0,
			StkCd:        stkCd,
			Market:       "KOSPI",
			Side:         "B",
			FilledQty:    qty,
			FilledPrice:  price,
			FilledAmount: qty * price,
			FeeAmount:    0.1,
			TaxAmount:    0.1,
		}

		tradeId, err := repository.InsertTradeLog(ctx, tx, virtualTradeLogEntity)
		if err != nil {
			fmt.Println("InsertTradeLog", err.Error())
		}
		fmt.Println(tradeId)

		// !!upsert 가상 자산 테이블 에
		virtualAssetEntity := model.TbVirtualAssetEntity{
			UserId:       0,
			AccountId:    0,
			StkCd:        stkCd,
			Market:       "KOSPI",
			PositionSide: "B",
			Qty:          qty,
			AvgPrice:     price,
			Status:       "ACTIVE",
		}
		err = repository.InsertVirtualAsset(ctx, tx, virtualAssetEntity)
		if err != nil {
			fmt.Println("InsertVirtualAsset", err.Error())
		}

		// !!update 가상 계좌 테이블 거래 가능 금액
		virtualAccountEntity := model.TbVirtualAccountEntity{
			AccountId: 0,
		}
		err = repository.UpdateVirtualAccount(ctx, tx, virtualAccountEntity, "BUY", int64(price*qty))
		// 트랜잭션으로 묶어서 commit
		if err := tx.Commit(ctx); err != nil {
			logger.Error("매수하다가 오류났다 오류.")
		}

	}
}
