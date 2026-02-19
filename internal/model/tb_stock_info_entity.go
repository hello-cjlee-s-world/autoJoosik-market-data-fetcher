package model

import (
	"autoJoosik-market-data-fetcher/pkg/logger"
	"encoding/json"
	"fmt"
	"strings"
)

// TbStockInfoEntity 주식 기본 정보
type TbStockInfoEntity struct {
	StkCd         string `db:"stk_cd" json:"stk_cd"`
	StkNm         string `db:"stk_nm" json:"stk_nm"`
	SetlMm        string `db:"setl_mm" json:"setl_mm"`
	Fav           string `db:"fav" json:"fav"`
	Cap           string `db:"cap" json:"cap"`
	FloStk        string `db:"flo_stk" json:"flo_stk"`
	CrdRt         string `db:"crd_rt" json:"crd_rt"`
	OyrHgst       string `db:"oyr_hgst" json:"oyr_hgst"`
	OyrLwst       string `db:"oyr_lwst" json:"oyr_lwst"`
	Mac           string `db:"mac" json:"mac"`
	MacWght       string `db:"mac_wght" json:"mac_wght"`
	ForExhRt      string `db:"for_exh_rt" json:"for_exh_rt"`
	ReplPric      string `db:"repl_pric" json:"repl_pric"`
	Per           string `db:"per" json:"per"`
	Eps           string `db:"eps" json:"eps"`
	Roe           string `db:"roe" json:"roe"`
	Pbr           string `db:"pbr" json:"pbr"`
	Ev            string `db:"ev" json:"ev"`
	Bps           string `db:"bps" json:"bps"`
	SaleAmt       string `db:"sale_amt" json:"sale_amt"`
	BusPro        string `db:"bus_pro" json:"bus_pro"`
	CupNga        string `db:"cup_nga" json:"cup_nga"`
	Hgst250       string `db:"250hgst" json:"250hgst"`
	Lwst250       string `db:"250lwst" json:"250lwst"`
	HighPric      string `db:"high_pric" json:"high_pric"`
	OpenPric      string `db:"open_pric" json:"open_pric"`
	LowPric       string `db:"low_pric" json:"low_pric"`
	UplPric       string `db:"upl_pric" json:"upl_pric"`
	LstPric       string `db:"lst_pric" json:"lst_pric"`
	BasePric      string `db:"base_pric" json:"base_pric"`
	ExpCntrPric   string `db:"exp_cntr_pric" json:"exp_cntr_pric"`
	ExpCntrQty    string `db:"exp_cntr_qty" json:"exp_cntr_qty"`
	Hgst250PricDt string `db:"250hgst_pric_dt" json:"250hgst_pric_dt"`
	Hgst250PreRt  string `db:"250hgst_pric_pre_rt" json:"250hgst_pric_pre_rt"`
	Lwst250PricDt string `db:"250lwst_pric_dt" json:"250lwst_pric_dt"`
	Lwst250PreRt  string `db:"250lwst_pric_pre_rt" json:"250lwst_pric_pre_rt"`
	CurPrc        string `db:"cur_prc" json:"cur_prc"`
	PreSig        string `db:"pre_sig" json:"pre_sig"`
	PredPre       string `db:"pred_pre" json:"pred_pre"`
	FluRt         string `db:"flu_rt" json:"flu_rt"`
	TrdeQty       string `db:"trde_qty" json:"trde_qty"`
	TrdePre       string `db:"trde_pre" json:"trde_pre"`
	FavUnit       string `db:"fav_unit" json:"fav_unit"`
	DstrStk       string `db:"dstr_stk" json:"dstr_stk"`
	DstrRt        string `db:"dstr_rt" json:"dstr_rt"`
	UpdatedAt     string `db:"updated_at" json:"updated_at"`
}

func ToTbStockInfoEntity(str string) (TbStockInfoEntity, error) {
	return ToTbStockInfoEntityWithFallback(str, "")
}

func ToTbStockInfoEntityWithFallback(str string, fallbackStkCd string) (TbStockInfoEntity, error) {
	var entity TbStockInfoEntity
	err := json.Unmarshal([]byte(str), &entity)
	if err != nil {
		logger.Error("While doing ToTbStockInfoEntity :: ", err.Error())
		return TbStockInfoEntity{}, err
	}

	entity.StkCd = strings.TrimSpace(entity.StkCd)
	if entity.StkCd == "" {
		fallbackStkCd = strings.TrimSpace(fallbackStkCd)
		if fallbackStkCd == "" {
			err = fmt.Errorf("empty stk_cd in stock info payload")
			logger.Error("While doing ToTbStockInfoEntity :: ", err.Error())
			return TbStockInfoEntity{}, err
		}
		logger.Warn("ToTbStockInfoEntity :: use fallback stk_cd", "stkCd", fallbackStkCd)
		entity.StkCd = fallbackStkCd
	}

	return entity, nil
}
