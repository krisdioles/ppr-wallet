package domain

import (
	"context"
	"time"
)

type UserBalance struct {
	ID          int64     `json:"id" db:"id"`
	Username    string    `json:"username" db:"username"`
	Balance     int64     `json:"balance" db:"balance"`
	BankCode    string    `json:"bank_code" db:"bank_code"`
	AccountNo   string    `json:"account_no" db:"account_no"`
	AccountName string    `json:"account_name" db:"account_name"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func (u *UserBalance) TableName() string {
	return "user_balances"
}

type UserBalanceRepository interface {
	GetByID(ctx context.Context, id int64) (*UserBalance, error)
	UpdateBalanceByID(ctx context.Context, updatedBalance, id int64) error
}

type UserBalanceUsecase interface {
	GetUserBalanceByID(ctx context.Context, id int64) (*UserBalance, error)
	DisburseBalance(ctx context.Context, id int64) error
}
