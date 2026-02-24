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
	"time"
)

func Buy(stkCd string, qty float64) error {
	rst, err := kiwoomApi.GetOrderBookLog(stkCd)
	if !utils.IsTradableTime(time.Now()) {
		return fmt.Errorf("not available trade time")
	}
	if err == nil {
		//  주식 거래 예시 트랜잭션으로 묶기
		ctx := context.Background()
		pool := datasource.GetPool()
		tx, err := pool.Begin(ctx)
		orderBookEntity := model.ToOrderBookLogEntity(rst)

		price := utils.ParseFloat(orderBookEntity.SelFprBid)
		//remainingQty, _ := strconv.ParseFloat(orderBookEntity.SelFprReq, 64)

		// !!insert 주문 정보(상태)
		//qty := 1.0          // 내가 사고 싶은 수량
		remainingQty := 0.0 // 내 주문 기준 남은 수량
		userId := int64(0)
		accountID := int64(0)
		clientOrderID := fmt.Sprintf("%d_%d_%d_%s",
			accountID,
			time.Now().UnixNano(),
			rand.Intn(1000),
			"buy",
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
			logger.Error("InsertOrderLog", err.Error())
			return err
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

		_, err = repository.InsertTradeLog(ctx, tx, virtualTradeLogEntity)
		if err != nil {
			logger.Error("InsertTradeLog", err.Error())
			return err
		}
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
		err = repository.UpsertVirtualAsset(ctx, tx, virtualAssetEntity)
		if err != nil {
			logger.Error("UpsertVirtualAsset", err.Error())
			return err
		}

		// !!update 가상 계좌 테이블 거래 가능 금액
		virtualAccountEntity := model.TbVirtualAccountEntity{
			AccountId: 0,
		}
		err = repository.UpdateVirtualAccount(ctx, tx, virtualAccountEntity, "BUY", int64(price*qty))
		if err != nil {
			logger.Error("UpdateVirtualAccount", err.Error())
			return err
		}
		// 트랜잭션으로 묶어서 commit
		if err := tx.Commit(ctx); err != nil {
			logger.Error("Buy :: 매수 도중 오류 발생")
			return err
		} else {
			logger.Info("Buy :: 매수 성공, stkCd=" + stkCd)
		}

	}

	return err
}
