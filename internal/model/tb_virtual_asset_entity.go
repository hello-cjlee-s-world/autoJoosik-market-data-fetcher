package model

import (
	"autoJoosik-market-data-fetcher/pkg/logger"
	"encoding/json"
	"time"
)

// StockInfoEntity 주식 기본 정보
type TbVirtualAssetEntity struct {
	AssetId        int64     `db:"asset_id" json:"asset_id"`
	UserId         int64     `db:"user_id" json:"user_id"`
	AccountId      int64     `db:"account_id" json:"account_id"`
	StkCd          string    `db:"stk_cd" json:"stk_cd"`
	Market         string    `db:"market" json:"market"`
	PositionSide   string    `db:"position_side" json:"position_side"`
	Qty            float64   `db:"qty" json:"qty"`
	AvailableQty   float64   `db:"available_qty" json:"available_qty"`
	AvgPrice       float64   `db:"avg_price" json:"avg_price"`
	LastPrice      float64   `db:"last_price" json:"last_price"`
	InvestedAmount float64   `db:"invested_amount" json:"invested_amount"`
	EvalAmount     float64   `db:"eval_amount" json:"eval_amount"`
	EvalPl         float64   `db:"eval_pl" json:"eval_pl"`
	EvalPlRate     float64   `db:"eval_pl_rate" json:"eval_pl_rate"`
	TodayBuyQty    float64   `db:"today_buy_qty" json:"today_buy_qty"`
	TodaySellQty   float64   `db:"today_sell_qty" json:"today_sell_qty"`
	Status         string    `db:"status" json:"status"`
	LastEvalAt     time.Time `db:"last_eval_at" json:"last_eval_at"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

func ToTbVirtualAssetEntity(str string) TbVirtualAssetEntity {
	var entity TbVirtualAssetEntity
	err := json.Unmarshal([]byte(str), &entity)
	if err != nil {
		logger.Error("While doing ToTbVirtualAssetEntity :: ", err.Error())
	}

	return entity
}
