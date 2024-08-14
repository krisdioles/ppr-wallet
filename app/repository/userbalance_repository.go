package repository

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/krisdioles/ppr-wallet/app/domain"
)

type UserBalanceRepository struct {
	DB *sqlx.DB
}

func NewUserBalanceRepository(db *sqlx.DB) *UserBalanceRepository {
	return &UserBalanceRepository{
		DB: db,
	}
}

func (r *UserBalanceRepository) GetByID(ctx context.Context, id int64) (*domain.UserBalance, error) {
	getByIDQuery := `SELECT * FROM user_balances WHERE id = ?`

	var userBalance = &domain.UserBalance{}
	if err := r.DB.GetContext(ctx, userBalance, getByIDQuery, id); err != nil {
		log.Println("[GetByID] query err:", err)
		return userBalance, err
	}

	return userBalance, nil
}

func (r *UserBalanceRepository) UpdateBalanceByID(ctx context.Context, updatedBalance, id int64) error {
	updateBalanceByIDQuery := `UPDATE user_balances SET balance = ? WHERE id = ?`

	_, err := r.DB.ExecContext(ctx, updateBalanceByIDQuery, updatedBalance, id)
	if err != nil {
		log.Println("[UpdateBalanceByID] query err:", err)
		return err
	}

	return nil
}
