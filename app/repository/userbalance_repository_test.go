package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/krisdioles/ppr-wallet/app/domain"
	"github.com/krisdioles/ppr-wallet/app/repository"
	"github.com/stretchr/testify/assert"
)

func TestUserBalanceRepository_GetByID(t *testing.T) {
	// Create a mock DB and expect the query
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.UserBalanceRepository{DB: sqlxDB}

	userID := int64(1)
	expectedUserBalance := &domain.UserBalance{
		ID:      userID,
		Balance: 1000,
	}

	rows := sqlmock.NewRows([]string{"id", "balance"}).
		AddRow(expectedUserBalance.ID, expectedUserBalance.Balance)

	mock.ExpectQuery("SELECT \\* FROM user_balances WHERE id = ?").
		WithArgs(userID).
		WillReturnRows(rows)

	// Execute the function
	result, err := repo.GetByID(context.Background(), userID)

	// Assert the expectations
	assert.NoError(t, err)
	assert.Equal(t, expectedUserBalance, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserBalanceRepository_GetByID_Error(t *testing.T) {
	// Create a mock DB and expect the query
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.UserBalanceRepository{DB: sqlxDB}

	userID := int64(1)

	mock.ExpectQuery("SELECT \\* FROM user_balances WHERE id = ?").
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	// Execute the function
	result, err := repo.GetByID(context.Background(), userID)

	// Assert the expectations
	assert.Error(t, err)
	assert.Equal(t, &domain.UserBalance{}, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserBalanceRepository_UpdateBalanceByID(t *testing.T) {
	// Create a mock DB and expect the exec
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.UserBalanceRepository{DB: sqlxDB}

	userID := int64(1)
	newBalance := int64(0)

	mock.ExpectExec("UPDATE user_balances SET balance = \\? WHERE id = \\?").
		WithArgs(newBalance, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Execute the function
	err = repo.UpdateBalanceByID(context.Background(), newBalance, userID)

	// Assert the expectations
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserBalanceRepository_UpdateBalanceByID_Error(t *testing.T) {
	// Create a mock DB and expect the exec
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.UserBalanceRepository{DB: sqlxDB}

	userID := int64(1)
	newBalance := int64(2000)

	mock.ExpectExec("UPDATE user_balances SET balance = \\? WHERE id = \\?").
		WithArgs(newBalance, userID).
		WillReturnError(errors.New("some error"))

	// Execute the function
	err = repo.UpdateBalanceByID(context.Background(), newBalance, userID)

	// Assert the expectations
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
