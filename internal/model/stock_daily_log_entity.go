package model

import (
	"autoJoosik-market-data-fetcher/pkg/logger"
	"encoding/json"
	"time"
)

// StockDailyLogEntity 주식 일/주/월/시/분 로그
type StockDailyLogEntity struct {
	ID    int64     `db:"id"`     // 로그 PK
	StkCd string    `db:"stk_cd"` // 종목코드
	Date  time.Time `db:"date"`   // 날짜

	OpenPric  float64 `db:"open_pric"`  // 시가
	HighPric  float64 `db:"high_pric"`  // 고가
	LowPric   float64 `db:"low_pric"`   // 저가
	ClosePric float64 `db:"close_pric"` // 종가
	Pre       float64 `db:"pre"`        // 대비
	FluRt     float64 `db:"flu_rt"`     // 등락률
	TrdeQty   int64   `db:"trde_qty"`   // 거래량
	TrdePrica float64 `db:"trde_prica"` // 거래대금

	ForPoss     int64   `db:"for_poss"`     // 외인보유
	ForWght     float64 `db:"for_wght"`     // 외인비중
	ForNetprps  int64   `db:"for_netprps"`  // 외인순매수
	OrgnNetprps int64   `db:"orgn_netprps"` // 기관순매수
	IndNetprps  int64   `db:"ind_netprps"`  // 개인순매수
	CrdRemnRt   float64 `db:"crd_remn_rt"`  // 신용잔고율
	Frgn        int64   `db:"frgn"`         // 외국계
	Prm         int64   `db:"prm"`          // 프로그램

	CreatedAt time.Time `db:"created_at"` // 적재 시각
}

func ToStockDailyLogEntity(str string) StockDailyLogEntity {
	var entity StockDailyLogEntity
	err := json.Unmarshal([]byte(str), &entity)
	if err != nil {
		logger.Error("While doing ToStockDailyLogEntity :: ", err.Error())
	}

	return entity
}
