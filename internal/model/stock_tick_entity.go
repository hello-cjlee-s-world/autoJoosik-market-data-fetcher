package model

import (
	"autoJoosik-market-data-fetcher/pkg/logger"
	"encoding/json"
	"time"
)

// StockTickEntity 주식 틱차트
type StockTickEntity struct {
	ID         int64     `db:"id"`           // 로그 PK
	StkCd      string    `db:"stk_cd"`       // 종목코드
	LastTicCnt int64     `db:"last_tic_cnt"` // 마지막 틱 카운트(공급원 정의에 따름)
	CurPrc     float64   `db:"cur_prc"`      // 현재가
	TrdeQty    int64     `db:"trde_qty"`     // 체결수량(거래량)
	CntrTm     time.Time `db:"cntr_tm"`      // 체결시각(타임존/포맷 주의)

	OpenPric float64 `db:"open_pric"` // 시가
	HighPric float64 `db:"high_pric"` // 고가
	LowPric  float64 `db:"low_pric"`  // 저가

	UpdStkpcTp string  `db:"upd_stkpc_tp"` // 등락 구분 코드(상승/하락/보합 등)
	UpdRt      float64 `db:"upd_rt"`       // 등락률(%)

	BicIndsTp string `db:"bic_inds_tp"` // 업종(대분류) 코드
	SmIndsTp  string `db:"sm_inds_tp"`  // 업종(소분류) 코드
	StkInfr   string `db:"stk_infr"`    // 종목 속성(ETF/ETN/우선주 등 코드)

	UpdStkpcEvent string  `db:"upd_stkpc_event"` // 가격변동 이벤트/사유(있다면)
	PredClosePric float64 `db:"pred_close_pric"` // 예상체결가/예상종가(공급원 정의)
}

func ToStockTickEntity(str string) StockTickEntity {
	var entity StockTickEntity
	err := json.Unmarshal([]byte(str), &entity)
	if err != nil {
		logger.Error("While doing ToStockTickEntity :: ", err.Error())
	}

	return entity
}
