package model

import (
	"autoJoosik-market-data-fetcher/pkg/logger"
	"encoding/json"
	"time"
)

// AccountProfitLogEntity 계좌 수익률 로그
type AccountProfitLogEntity struct {
	ID    int64     `db:"id"`     // 로그 PK
	Dt    time.Time `db:"dt"`     // 일자
	StkCd string    `db:"stk_cd"` // 종목코드
	StkNm string    `db:"stk_nm"` // 종목명

	CurPrc  float64 `db:"cur_prc"`  // 현재가
	PurPric float64 `db:"pur_pric"` // 매입가
	PurAmt  float64 `db:"pur_amt"`  // 매입금액
	RmndQty int64   `db:"rmnd_qty"` // 보유수량

	TdySelPl    float64 `db:"tdy_sel_pl"`    // 당일매도손익
	TdyTrdeCmsn float64 `db:"tdy_trde_cmsn"` // 당일매매수수료
	TdyTrdeTax  float64 `db:"tdy_trde_tax"`  // 당일매매세금

	CrdTp       string    `db:"crd_tp"`        // 신용구분
	LoanDt      time.Time `db:"loan_dt"`       // 대출일
	SetlRemn    float64   `db:"setl_remn"`     // 결제잔고
	ClrnAlowQty int64     `db:"clrn_alow_qty"` // 청산가능수량
	CrdAmt      float64   `db:"crd_amt"`       // 신용금액
	CrdInt      float64   `db:"crd_int"`       // 신용이자
	ExprDt      time.Time `db:"expr_dt"`       // 만기일

	CreatedAt time.Time `db:"created_at"` // 적재 시각
}

func ToAccountProfitLogEntity(str string) AccountProfitLogEntity {
	var entity AccountProfitLogEntity
	err := json.Unmarshal([]byte(str), &entity)
	if err != nil {
		logger.Error("While doing ToAccountProfitLogEntity :: ", err.Error())
	}

	return entity
}
