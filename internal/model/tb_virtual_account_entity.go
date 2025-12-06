package model

import (
	"autoJoosik-market-data-fetcher/pkg/logger"
	"encoding/json"
	"time"
)

// TbVirtualAccountEntity 가상 계좌 테이블
type TbVirtualAccountEntity struct {
	AccountId      int64     `db:"account_id" json:"account_id"`
	UserId         int64     `db:"user_id" json:"user_id"`
	AccountName    string    `db:"account_name" json:"account_name"`
	CashBalance    float64   `db:"cash_balance" json:"cash_balance"`
	TotalInvested  float64   `db:"total_invested" json:"total_invested"`
	TotalEval      float64   `db:"total_eval" json:"total_eval"`
	TotalPl        float64   `db:"total_pl" json:"total_pl"`
	TotalPlRate    float64   `db:"total_pl_rate" json:"total_pl_rate"`
	DepositAmount  float64   `db:"deposit_amount" json:"deposit_amount"`
	WithdrawAmount float64   `db:"withdraw_amount" json:"withdraw_amount"`
	Status         string    `db:"status" json:"status"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

func ToTbVirtualAccountEntity(str string) TbVirtualAccountEntity {
	var entity TbVirtualAccountEntity
	err := json.Unmarshal([]byte(str), &entity)
	if err != nil {
		logger.Error("While doing ToTbVirtualAccountEntity :: ", err.Error())
	}

	return entity
}
