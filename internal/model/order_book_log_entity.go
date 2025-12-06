package model

import (
	"autoJoosik-market-data-fetcher/pkg/logger"
	"encoding/json"
	"time"
)

// OrderBookLogEntity 주식호가 로그
type OrderBookLogEntity struct {
	ID           string `db:"id" json:"id"`                           // 내부 PK
	BidReqBaseTm string `db:"bid_req_base_tm" json:"bid_req_base_tm"` // 호가잔량기준시간

	// 매도호가 (10차~1차)
	Sel10thPreReqPre string `db:"sel_10th_pre_req_pre" json:"sel_10th_pre_req_pre"`
	Sel10thPreReq    string `db:"sel_10th_pre_req" json:"sel_10th_pre_req"`
	Sel10thPreBid    string `db:"sel_10th_pre_bid" json:"sel_10th_pre_bid"`
	Sel9thPreReqPre  string `db:"sel_9th_pre_req_pre" json:"sel_9th_pre_req_pre"`
	Sel9thPreReq     string `db:"sel_9th_pre_req" json:"sel_9th_pre_req"`
	Sel9thPreBid     string `db:"sel_9th_pre_bid" json:"sel_9th_pre_bid"`
	Sel8thPreReqPre  string `db:"sel_8th_pre_req_pre" json:"sel_8th_pre_req_pre"`
	Sel8thPreReq     string `db:"sel_8th_pre_req" json:"sel_8th_pre_req"`
	Sel8thPreBid     string `db:"sel_8th_pre_bid" json:"sel_8th_pre_bid"`
	Sel7thPreReqPre  string `db:"sel_7th_pre_req_pre" json:"sel_7th_pre_req_pre"`
	Sel7thPreReq     string `db:"sel_7th_pre_req" json:"sel_7th_pre_req"`
	Sel7thPreBid     string `db:"sel_7th_pre_bid" json:"sel_7th_pre_bid"`
	Sel6thPreReqPre  string `db:"sel_6th_pre_req_pre" json:"sel_6th_pre_req_pre"`
	Sel6thPreReq     string `db:"sel_6th_pre_req" json:"sel_6th_pre_req"`
	Sel6thPreBid     string `db:"sel_6th_pre_bid" json:"sel_6th_pre_bid"`
	Sel5thPreReqPre  string `db:"sel_5th_pre_req_pre" json:"sel_5th_pre_req_pre"`
	Sel5thPreReq     string `db:"sel_5th_pre_req" json:"sel_5th_pre_req"`
	Sel5thPreBid     string `db:"sel_5th_pre_bid" json:"sel_5th_pre_bid"`
	Sel4thPreReqPre  string `db:"sel_4th_pre_req_pre" json:"sel_4th_pre_req_pre"`
	Sel4thPreReq     string `db:"sel_4th_pre_req" json:"sel_4th_pre_req"`
	Sel4thPreBid     string `db:"sel_4th_pre_bid" json:"sel_4th_pre_bid"`
	Sel3thPreReqPre  string `db:"sel_3th_pre_req_pre" json:"sel_3th_pre_req_pre"`
	Sel3thPreReq     string `db:"sel_3th_pre_req" json:"sel_3th_pre_req"`
	Sel3thPreBid     string `db:"sel_3th_pre_bid" json:"sel_3th_pre_bid"`
	Sel2thPreReqPre  string `db:"sel_2th_pre_req_pre" json:"sel_2th_pre_req_pre"`
	Sel2thPreReq     string `db:"sel_2th_pre_req" json:"sel_2th_pre_req"`
	Sel2thPreBid     string `db:"sel_2th_pre_bid" json:"sel_2th_pre_bid"`
	Sel1thPreReqPre  string `db:"sel_1th_pre_req_pre" json:"sel_1th_pre_req_pre"`
	SelFprReq        string `db:"sel_fpr_req" json:"sel_fpr_req"` // 매도최우선잔량
	SelFprBid        string `db:"sel_fpr_bid" json:"sel_fpr_bid"` // 매도최우선호가

	// 매수호가 (1차~10차)
	BuyFprBid        string `db:"buy_fpr_bid" json:"buy_fpr_bid"` // 매수최우선호가
	BuyFprReq        string `db:"buy_fpr_req" json:"buy_fpr_req"` // 매수최우선잔량
	Buy1thPreReqPre  string `db:"buy_1th_pre_req_pre" json:"buy_1th_pre_req_pre"`
	Buy2thPreBid     string `db:"buy_2th_pre_bid" json:"buy_2th_pre_bid"`
	Buy2thPreReq     string `db:"buy_2th_pre_req" json:"buy_2th_pre_req"`
	Buy2thPreReqPre  string `db:"buy_2th_pre_req_pre" json:"buy_2th_pre_req_pre"`
	Buy3thPreBid     string `db:"buy_3th_pre_bid" json:"buy_3th_pre_bid"`
	Buy3thPreReq     string `db:"buy_3th_pre_req" json:"buy_3th_pre_req"`
	Buy3thPreReqPre  string `db:"buy_3th_pre_req_pre" json:"buy_3th_pre_req_pre"`
	Buy4thPreBid     string `db:"buy_4th_pre_bid" json:"buy_4th_pre_bid"`
	Buy4thPreReq     string `db:"buy_4th_pre_req" json:"buy_4th_pre_req"`
	Buy4thPreReqPre  string `db:"buy_4th_pre_req_pre" json:"buy_4th_pre_req_pre"`
	Buy5thPreBid     string `db:"buy_5th_pre_bid" json:"buy_5th_pre_bid"`
	Buy5thPreReq     string `db:"buy_5th_pre_req" json:"buy_5th_pre_req"`
	Buy5thPreReqPre  string `db:"buy_5th_pre_req_pre" json:"buy_5th_pre_req_pre"`
	Buy6thPreBid     string `db:"buy_6th_pre_bid" json:"buy_6th_pre_bid"`
	Buy6thPreReq     string `db:"buy_6th_pre_req" json:"buy_6th_pre_req"`
	Buy6thPreReqPre  string `db:"buy_6th_pre_req_pre" json:"buy_6th_pre_req_pre"`
	Buy7thPreBid     string `db:"buy_7th_pre_bid" json:"buy_7th_pre_bid"`
	Buy7thPreReq     string `db:"buy_7th_pre_req" json:"buy_7th_pre_req"`
	Buy7thPreReqPre  string `db:"buy_7th_pre_req_pre" json:"buy_7th_pre_req_pre"`
	Buy8thPreBid     string `db:"buy_8th_pre_bid" json:"buy_8th_pre_bid"`
	Buy8thPreReq     string `db:"buy_8th_pre_req" json:"buy_8th_pre_req"`
	Buy8thPreReqPre  string `db:"buy_8th_pre_req_pre" json:"buy_8th_pre_req_pre"`
	Buy9thPreBid     string `db:"buy_9th_pre_bid" json:"buy_9th_pre_bid"`
	Buy9thPreReq     string `db:"buy_9th_pre_req" json:"buy_9th_pre_req"`
	Buy9thPreReqPre  string `db:"buy_9th_pre_req_pre" json:"buy_9th_pre_req_pre"`
	Buy10thPreBid    string `db:"buy_10th_pre_bid" json:"buy_10th_pre_bid"`
	Buy10thPreReq    string `db:"buy_10th_pre_req" json:"buy_10th_pre_req"`
	Buy10thPreReqPre string `db:"buy_10th_pre_req_pre" json:"buy_10th_pre_req_pre"`

	// 총잔량
	TotSelReqJubPre string `db:"tot_sel_req_jub_pre" json:"tot_sel_req_jub_pre"`
	TotSelReq       string `db:"tot_sel_req" json:"tot_sel_req"`
	TotBuyReq       string `db:"tot_buy_req" json:"tot_buy_req"`
	TotBuyReqJubPre string `db:"tot_buy_req_jub_pre" json:"tot_buy_req_jub_pre"`

	// 시간외 잔량
	OvtSelReqPre string `db:"ovt_sel_req_pre" json:"ovt_sel_req_pre"`
	OvtSelReq    string `db:"ovt_sel_req" json:"ovt_sel_req"`
	OvtBuyReq    string `db:"ovt_buy_req" json:"ovt_buy_req"`
	OvtBuyReqPre string `db:"ovt_buy_req_pre" json:"ovt_buy_req_pre"`

	CreatedAt time.Time `db:"created_at" json:"created_at"` // 로그 적재 시각
}

func ToOrderBookLogEntity(str string) OrderBookLogEntity {
	var entity OrderBookLogEntity
	err := json.Unmarshal([]byte(str), &entity)
	if err != nil {
		logger.Error("While doing ToOrderBookLogEntity :: ", err.Error())
	}

	return entity
}
