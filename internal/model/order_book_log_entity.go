package model

import (
	"autoJoosik-market-data-fetcher/pkg/logger"
	"encoding/json"
	"time"
)

// OrderBookLogEntity 주식호가 로그
type OrderBookLogEntity struct {
	ID           int64  `db:"id"`              // 내부 PK
	BidReqBaseTm string `db:"bid_req_base_tm"` // 호가잔량기준시간

	// 매도호가 (10차~1차)
	Sel10thPreReqPre float64 `db:"sel_10th_pre_req_pre"`
	Sel10thPreReq    int64   `db:"sel_10th_pre_req"`
	Sel10thPreBid    float64 `db:"sel_10th_pre_bid"`
	Sel9thPreReqPre  float64 `db:"sel_9th_pre_req_pre"`
	Sel9thPreReq     int64   `db:"sel_9th_pre_req"`
	Sel9thPreBid     float64 `db:"sel_9th_pre_bid"`
	Sel8thPreReqPre  float64 `db:"sel_8th_pre_req_pre"`
	Sel8thPreReq     int64   `db:"sel_8th_pre_req"`
	Sel8thPreBid     float64 `db:"sel_8th_pre_bid"`
	Sel7thPreReqPre  float64 `db:"sel_7th_pre_req_pre"`
	Sel7thPreReq     int64   `db:"sel_7th_pre_req"`
	Sel7thPreBid     float64 `db:"sel_7th_pre_bid"`
	Sel6thPreReqPre  float64 `db:"sel_6th_pre_req_pre"`
	Sel6thPreReq     int64   `db:"sel_6th_pre_req"`
	Sel6thPreBid     float64 `db:"sel_6th_pre_bid"`
	Sel5thPreReqPre  float64 `db:"sel_5th_pre_req_pre"`
	Sel5thPreReq     int64   `db:"sel_5th_pre_req"`
	Sel5thPreBid     float64 `db:"sel_5th_pre_bid"`
	Sel4thPreReqPre  float64 `db:"sel_4th_pre_req_pre"`
	Sel4thPreReq     int64   `db:"sel_4th_pre_req"`
	Sel4thPreBid     float64 `db:"sel_4th_pre_bid"`
	Sel3thPreReqPre  float64 `db:"sel_3th_pre_req_pre"`
	Sel3thPreReq     int64   `db:"sel_3th_pre_req"`
	Sel3thPreBid     float64 `db:"sel_3th_pre_bid"`
	Sel2thPreReqPre  float64 `db:"sel_2th_pre_req_pre"`
	Sel2thPreReq     int64   `db:"sel_2th_pre_req"`
	Sel2thPreBid     float64 `db:"sel_2th_pre_bid"`
	Sel1thPreReqPre  float64 `db:"sel_1th_pre_req_pre"`
	SelFprReq        int64   `db:"sel_fpr_req"` // 매도최우선잔량
	SelFprBid        float64 `db:"sel_fpr_bid"` // 매도최우선호가

	// 매수호가 (1차~10차)
	BuyFprBid        float64 `db:"buy_fpr_bid"` // 매수최우선호가
	BuyFprReq        int64   `db:"buy_fpr_req"` // 매수최우선잔량
	Buy1thPreReqPre  float64 `db:"buy_1th_pre_req_pre"`
	Buy2thPreBid     float64 `db:"buy_2th_pre_bid"`
	Buy2thPreReq     int64   `db:"buy_2th_pre_req"`
	Buy2thPreReqPre  float64 `db:"buy_2th_pre_req_pre"`
	Buy3thPreBid     float64 `db:"buy_3th_pre_bid"`
	Buy3thPreReq     int64   `db:"buy_3th_pre_req"`
	Buy3thPreReqPre  float64 `db:"buy_3th_pre_req_pre"`
	Buy4thPreBid     float64 `db:"buy_4th_pre_bid"`
	Buy4thPreReq     int64   `db:"buy_4th_pre_req"`
	Buy4thPreReqPre  float64 `db:"buy_4th_pre_req_pre"`
	Buy5thPreBid     float64 `db:"buy_5th_pre_bid"`
	Buy5thPreReq     int64   `db:"buy_5th_pre_req"`
	Buy5thPreReqPre  float64 `db:"buy_5th_pre_req_pre"`
	Buy6thPreBid     float64 `db:"buy_6th_pre_bid"`
	Buy6thPreReq     int64   `db:"buy_6th_pre_req"`
	Buy6thPreReqPre  float64 `db:"buy_6th_pre_req_pre"`
	Buy7thPreBid     float64 `db:"buy_7th_pre_bid"`
	Buy7thPreReq     int64   `db:"buy_7th_pre_req"`
	Buy7thPreReqPre  float64 `db:"buy_7th_pre_req_pre"`
	Buy8thPreBid     float64 `db:"buy_8th_pre_bid"`
	Buy8thPreReq     int64   `db:"buy_8th_pre_req"`
	Buy8thPreReqPre  float64 `db:"buy_8th_pre_req_pre"`
	Buy9thPreBid     float64 `db:"buy_9th_pre_bid"`
	Buy9thPreReq     int64   `db:"buy_9th_pre_req"`
	Buy9thPreReqPre  float64 `db:"buy_9th_pre_req_pre"`
	Buy10thPreBid    float64 `db:"buy_10th_pre_bid"`
	Buy10thPreReq    int64   `db:"buy_10th_pre_req"`
	Buy10thPreReqPre float64 `db:"buy_10th_pre_req_pre"`

	// 총잔량
	TotSelReqJubPre float64 `db:"tot_sel_req_jub_pre"`
	TotSelReq       int64   `db:"tot_sel_req"`
	TotBuyReq       int64   `db:"tot_buy_req"`
	TotBuyReqJubPre float64 `db:"tot_buy_req_jub_pre"`

	// 시간외 잔량
	OvtSelReqPre float64 `db:"ovt_sel_req_pre"`
	OvtSelReq    int64   `db:"ovt_sel_req"`
	OvtBuyReq    int64   `db:"ovt_buy_req"`
	OvtBuyReqPre float64 `db:"ovt_buy_req_pre"`

	CreatedAt time.Time `db:"created_at"` // 로그 적재 시각
}

func ToOrderBookLogEntity(str string) OrderBookLogEntity {
	var entity OrderBookLogEntity
	err := json.Unmarshal([]byte(str), &entity)
	if err != nil {
		logger.Error("While doing ToOrderBookLogEntity :: ", err.Error())
	}

	return entity
}
