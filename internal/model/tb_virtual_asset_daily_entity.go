package model

import (
	"autoJoosik-market-data-fetcher/pkg/logger"
	"encoding/json"
	"time"
)

// TbVirtualAssetDailyEntity 가상자산 이력 테이블
type TbVirtualAssetDailyEntity struct {
	Id          int64     `db:"id" json:"id"`
	AssetId     int64     `db:"asset_id" json:"asset_id"`
	UserId      int64     `db:"user_id" json:"user_id"`
	AccountId   int64     `db:"account_id" json:"account_id"`
	BaseDate    string    `db:"base_date" json:"base_date"`
	TotalAssets string    `db:"total_assets" json:"total_assets"`
	StockValue  string    `db:"stock_value" json:"stock_value"`
	CashBalance float64   `db:"cash_balance" json:"cash_balance"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

func ToTbVirtualAssetDailyEntity(str string) TbVirtualAssetDailyEntity {
	var entity TbVirtualAssetDailyEntity
	err := json.Unmarshal([]byte(str), &entity)
	if err != nil {
		logger.Error("While doing ToTbVirtualAssetDailyEntity :: ", err.Error())
	}

	return entity
}
