package model

import (
	"autoJoosik-market-data-fetcher/pkg/logger"
	"encoding/json"
	"time"
)

// TbVirtualOrder 가상 주문 테이블
type TbVirtualOrder struct {
	OrderID       int64     `db:"order_id" json:"order_id"`
	UserID        int64     `db:"user_id" json:"user_id"`
	AccountID     int64     `db:"account_id" json:"account_id"`
	StkCd         string    `db:"stk_cd" json:"stk_cd"`
	Market        string    `db:"market" json:"market"`
	Side          string    `db:"side" json:"side"`
	OrderType     string    `db:"order_type" json:"order_type"`
	TimeInForce   string    `db:"time_in_force" json:"time_in_force"`
	Price         float64   `db:"price" json:"price"`
	Qty           float64   `db:"qty" json:"qty"`
	FilledQty     float64   `db:"filled_qty" json:"filled_qty"`
	RemainingQty  float64   `db:"remaining_qty" json:"remaining_qty"`
	Status        string    `db:"status" json:"status"`
	ClientOrderID string    `db:"client_order_id" json:"client_order_id"`
	Reason        string    `db:"reason" json:"reason"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

func ToTbVirtualOrder(str string) TbVirtualOrder {
	var entity TbVirtualOrder
	err := json.Unmarshal([]byte(str), &entity)
	if err != nil {
		logger.Error("While doing ToTbVirtualOrder :: ", err.Error())
	}

	return entity
}
