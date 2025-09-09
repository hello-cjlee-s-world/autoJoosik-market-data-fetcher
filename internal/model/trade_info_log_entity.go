package model

import (
	"autoJoosik-market-data-fetcher/internal/utils"
	"autoJoosik-market-data-fetcher/pkg/logger"
	"encoding/json"
	"time"
)

type TradeResponse struct {
	CntrInfr   []TradeInfoLogRaw `json:"cntr_infr"`
	ReturnCode int               `json:"return_code"`
	ReturnMsg  string            `json:"return_msg"`
}

// TradeInfoLogEntity 체결정보 로그
type TradeInfoLogRaw struct {
	Tm            string `json:"tm"`               // 시간 (API 그대로 저장)
	CurPrc        string `json:"cur_prc"`          // 현재가
	PredPre       string `json:"pred_pre"`         // 전일대비
	PreRt         string `json:"pre_rt"`           // 대비율 (%)
	PriSelBidUnit string `json:"pri_sel_bid_unit"` // 우선매도호가단위
	PriBuyBidUnit string `json:"pri_buy_bid_unit"` // 우선매수호가단위
	CntrTrdeQty   string `json:"cntr_trde_qty"`    // 체결거래량
	Sign          string `json:"sign"`             // sign
	AccTrdeQty    string `json:"acc_trde_qty"`     // 누적거래량
	AccTrdePrica  string `json:"acc_trde_prica"`   // 누적거래대금
	CntrStr       string `json:"cntr_str"`         // 체결강도
	StexTp        string `json:"stex_tp"`          // 거래소구분 (KRX, NXT, 통합)
}

type TradeInfoLogEntity struct {
	Tm            time.Time `db:"tm"`       // 시간
	CurPrc        float64   `db:"cur_prc"`  // 현재가
	PredPre       float64   `db:"pred_pre"` // 전일대비
	PreRt         float64   `db:"pre_rt"`   // 대비율
	PriSelBidUnit float64   `db:"pri_sel_bid_unit"`
	PriBuyBidUnit float64   `db:"pri_buy_bid_unit"`
	CntrTrdeQty   int64     `db:"cntr_trde_qty"`  // 체결거래량
	Sign          string    `db:"sign"`           // sign
	AccTrdeQty    int64     `db:"acc_trde_qty"`   // 누적거래량
	AccTrdePrica  float64   `db:"acc_trde_prica"` // 누적거래대금
	CntrStr       float64   `db:"cntr_str"`       // 체결강도
	StexTp        string    `db:"stex_tp"`        // 거래소구분
	StkCd         string    `db:"stk_cd"`         // 종목코드
	CreatedAt     time.Time `db:"created_at"`     // 적재 시간
}

func (r TradeInfoLogRaw) ToEntity(stkCd string) TradeInfoLogEntity {
	return TradeInfoLogEntity{
		Tm:            utils.ParseTradeTimestamp(r.Tm),
		CurPrc:        utils.ParseFloat(r.CurPrc),
		PredPre:       utils.ParseFloat(r.PredPre),
		PreRt:         utils.ParseFloat(r.PreRt),
		PriSelBidUnit: utils.ParseFloat(r.PriSelBidUnit),
		PriBuyBidUnit: utils.ParseFloat(r.PriBuyBidUnit),
		CntrTrdeQty:   utils.ParseInt(r.CntrTrdeQty),
		Sign:          r.Sign,
		AccTrdeQty:    utils.ParseInt(r.AccTrdeQty),
		AccTrdePrica:  utils.ParseFloat(r.AccTrdePrica),
		CntrStr:       utils.ParseFloat(r.CntrStr),
		StexTp:        r.StexTp,
		CreatedAt:     time.Now(),
		StkCd:         stkCd,
	}
}

func ToTradeInfoLogEntity(str string, stkCd string) []TradeInfoLogEntity {
	var resp TradeResponse
	err := json.Unmarshal([]byte(str), &resp)
	if err != nil {
		logger.Error("While doing ToTradeInfoLogEntity :: ", err.Error())
		return nil
	}

	entities := make([]TradeInfoLogEntity, 0, len(resp.CntrInfr))
	for _, raw := range resp.CntrInfr {
		entities = append(entities, raw.ToEntity(stkCd))
	}

	return entities
}
