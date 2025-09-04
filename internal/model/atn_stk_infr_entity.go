package model

import (
	"autoJoosik-market-data-fetcher/pkg/logger"
	"encoding/json"
)

// AtnStkInfrEntity 관심종목 정보
type AtnStkInfrEntity struct {
	StkCd          string `db:"stk_cd"`            // 종목코드
	StkNm          string `db:"stk_nm"`            // 종목명
	CurPrc         string `db:"cur_prc"`           // 현재가
	BasePric       string `db:"base_pric"`         // 기준가
	PredPre        string `db:"pred_pre"`          // 전일대비
	PredPreSig     string `db:"pred_pre_sig"`      // 전일대비기호
	FluRt          string `db:"flu_rt"`            // 등락율
	TrdeQty        string `db:"trde_qty"`          // 거래량
	TrdePrica      string `db:"trde_prica"`        // 거래대금
	CntrQty        string `db:"cntr_qty"`          // 체결량
	CntrStr        string `db:"cntr_str"`          // 체결강도
	PredTrdeQtyPre string `db:"pred_trde_qty_pre"` // 전일거래량대비
	SelBid         string `db:"sel_bid"`           // 매도호가
	BuyBid         string `db:"buy_bid"`           // 매수호가
	Sel1thBid      string `db:"sel_1th_bid"`
	Sel2thBid      string `db:"sel_2th_bid"`
	Sel3thBid      string `db:"sel_3th_bid"`
	Sel4thBid      string `db:"sel_4th_bid"`
	Sel5thBid      string `db:"sel_5th_bid"`
	Buy1thBid      string `db:"buy_1th_bid"`
	Buy2thBid      string `db:"buy_2th_bid"`
	Buy3thBid      string `db:"buy_3th_bid"`
	Buy4thBid      string `db:"buy_4th_bid"`
	Buy5thBid      string `db:"buy_5th_bid"`
	UplPric        string `db:"upl_pric"`      // 상한가
	LstPric        string `db:"lst_pric"`      // 하한가
	OpenPric       string `db:"open_pric"`     // 시가
	HighPric       string `db:"high_pric"`     // 고가
	LowPric        string `db:"low_pric"`      // 저가
	ClosePric      string `db:"close_pric"`    // 종가
	CntrTm         string `db:"cntr_tm"`       // 체결시간
	ExpCntrPric    string `db:"exp_cntr_pric"` // 예상체결가
	ExpCntrQty     string `db:"exp_cntr_qty"`  // 예상체결량
	Cap            string `db:"cap"`           // 자본금
	Fav            string `db:"fav"`           // 액면가
	Mac            string `db:"mac"`           // 시가총액
	Stkcnt         string `db:"stkcnt"`        // 주식수
	BidTm          string `db:"bid_tm"`        // 호가시간
	Dt             string `db:"dt"`            // 일자
	PriSelReq      string `db:"pri_sel_req"`   // 우선매도잔량
	PriBuyReq      string `db:"pri_buy_req"`   // 우선매수잔량
	PriSelCnt      string `db:"pri_sel_cnt"`   // 우선매도건수
	PriBuyCnt      string `db:"pri_buy_cnt"`   // 우선매수건수
	TotSelReq      string `db:"tot_sel_req"`   // 총매도잔량
	TotBuyReq      string `db:"tot_buy_req"`   // 총매수잔량
	TotSelCnt      string `db:"tot_sel_cnt"`   // 총매도건수
	TotBuyCnt      string `db:"tot_buy_cnt"`   // 총매수건수
	Prty           string `db:"prty"`          // 패리티
	Gear           string `db:"gear"`          // 기어링
	PlQutr         string `db:"pl_qutr"`       // 손익분기
	CapSupport     string `db:"cap_support"`   // 자본지지
	ElwexecPric    string `db:"elwexec_pric"`  // ELW행사가
	CnvtRt         string `db:"cnvt_rt"`       // 전환비율
	ElwexprDt      string `db:"elwexpr_dt"`    // ELW만기일
	CntrEngg       string `db:"cntr_engg"`     // 미결제약정
	CntrPredPre    string `db:"cntr_pred_pre"` // 미결제전일대비
	TheoryPric     string `db:"theory_pric"`   // 이론가
	InnrVltl       string `db:"innr_vltl"`     // 내재변동성
	Delta          string `db:"delta"`         // 델타
	Gam            string `db:"gam"`           // 감마
	Theta          string `db:"theta"`         // 쎄타
	Vega           string `db:"vega"`          // 베가
	Law            string `db:"law"`           // 로
}

func ToAtnStkInfrEntity(str string) AtnStkInfrEntity {
	var entity AtnStkInfrEntity
	err := json.Unmarshal([]byte(str), &entity)
	if err != nil {
		logger.Error("While doing ToAccountProfitLogEntity :: ", err.Error())
	}

	return entity
}
