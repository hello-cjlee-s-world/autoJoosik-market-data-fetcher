package repository

import (
	"autoJoosik-market-data-fetcher/internal/model"
	"autoJoosik-market-data-fetcher/pkg/logger"
	"context"
	"fmt"
)

func UpdateVirtualAccount(ctx context.Context, db DB, entity model.TbVirtualAccountEntity, flag string, money int64) error {
	// money: 체결금액 + 수수료 포함 (BUY는 실제로 빠지는 금액, SELL은 실제로 들어오는 금액)
	if money < 0 {
		return fmt.Errorf("money must be >= 0")
	}

	// 1) 현금/입출금/상태 업데이트
	updateQ, args, err := buildAccountCashflowSQL(flag, money, entity.AccountId)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, updateQ, args...)
	if err != nil {
		logger.Error("UpdateVirtualAccount :: cashflow update error :: ", err)
		return err
	}

	// 2) 보유 포지션 합산으로 요약(total_invested/total_eval/손익/수익률) 재계산
	//
	// tb_virtual_asset 가정:
	// - invested_amount : 매수원가(수수료 포함) 누적
	// - eval_amount     : 현재 평가금액(수량*현재가)
	//
	// ※ 너 테이블 컬럼명이 다르면 여기만 교체하면 됨.
	recalcQ := `
WITH s AS (
  SELECT
    COALESCE(SUM(invested_amount), 0) AS invested_sum,
    COALESCE(SUM(eval_amount), 0)     AS eval_sum
  FROM tb_virtual_asset
  WHERE account_id = $1
)
UPDATE tb_virtual_account a
SET
  total_invested = (SELECT invested_sum FROM s),
  total_eval     = (SELECT eval_sum FROM s) + a.cash_balance,
  total_pl       = ((SELECT eval_sum FROM s) + a.cash_balance) - (a.deposit_amount - a.withdraw_amount),
  total_pl_rate  = CASE
                     WHEN (a.deposit_amount - a.withdraw_amount) = 0 THEN 0
                     ELSE (
                       (((SELECT eval_sum FROM s) + a.cash_balance) - (a.deposit_amount - a.withdraw_amount))::numeric
                       / (a.deposit_amount - a.withdraw_amount)::numeric
                     ) * 100
                   END,
  updated_at     = now()
WHERE a.account_id = $1
`
	_, err = db.Exec(ctx, recalcQ, entity.AccountId)
	if err != nil {
		logger.Error("UpdateVirtualAccount :: recalc error :: ", err)
		return err
	}

	logger.Debug("UpdateVirtualAccount :: success :: ", "flag", flag, "money", money, "account_id", entity.AccountId)
	return nil
}

func buildAccountCashflowSQL(flag string, money int64, accountId int64) (string, []any, error) {
	switch flag {
	case "BUY":
		// 매수: 현금 감소 (수수료 포함)
		q := `
UPDATE tb_virtual_account
SET cash_balance = cash_balance - $1,
    updated_at   = now()
WHERE account_id = $2
`
		return q, []any{money, accountId}, nil

	case "SELL":
		// 매도: 현금 증가 (수수료 포함해서 실제 유입 금액)
		q := `
UPDATE tb_virtual_account
SET cash_balance = cash_balance + $1,
    updated_at   = now()
WHERE account_id = $2
`
		return q, []any{money, accountId}, nil

	case "DEPOSIT":
		q := `
UPDATE tb_virtual_account
SET cash_balance   = cash_balance + $1,
    deposit_amount = deposit_amount + $1,
    updated_at     = now()
WHERE account_id = $2
`
		return q, []any{money, accountId}, nil

	case "WITHDRAW":
		q := `
UPDATE tb_virtual_account
SET cash_balance    = cash_balance - $1,
    withdraw_amount = withdraw_amount + $1,
    updated_at      = now()
WHERE account_id = $2
`
		return q, []any{money, accountId}, nil

	default:
		return "", nil, fmt.Errorf("invalid flag=%s", flag)
	}
}
