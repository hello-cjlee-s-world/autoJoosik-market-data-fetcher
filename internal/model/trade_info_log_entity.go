package model

import (
	"autoJoosik-market-data-fetcher/pkg/logger"
	"encoding/json"
	"time"
)

// TradeInfoLogEntity 체결정보 로그
type TradeInfoLogEntity struct {
	Tm            string    `db:"tm"`               // 시간 (API 그대로 저장)
	CurPrc        float64   `db:"cur_prc"`          // 현재가
	PredPre       float64   `db:"pred_pre"`         // 전일대비
	PreRt         float64   `db:"pre_rt"`           // 대비율 (%)
	PriSelBidUnit float64   `db:"pri_sel_bid_unit"` // 우선매도호가단위
	PriBuyBidUnit float64   `db:"pri_buy_bid_unit"` // 우선매수호가단위
	CntrTrdeQty   int64     `db:"cntr_trde_qty"`    // 체결거래량
	Sign          string    `db:"sign"`             // sign
	AccTrdeQty    int64     `db:"acc_trde_qty"`     // 누적거래량
	AccTrdePrica  float64   `db:"acc_trde_prica"`   // 누적거래대금
	CntrStr       float64   `db:"cntr_str"`         // 체결강도
	StexTp        string    `db:"stex_tp"`          // 거래소구분 (KRX, NXT, 통합)
	CreatedAt     time.Time `db:"created_at"`       // 로그 적재 시각
}

func ToTradeInfoLogEntity(str string) TradeInfoLogEntity {
	var entity TradeInfoLogEntity
	err := json.Unmarshal([]byte(str), &entity)
	if err != nil {
		logger.Error("While doing ToTradeInfoLogEntity :: ", err.Error())
	}

	return entity
}
