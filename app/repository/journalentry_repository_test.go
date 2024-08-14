package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/krisdioles/ppr-wallet/app/domain"
	"github.com/krisdioles/ppr-wallet/app/repository"
	"github.com/stretchr/testify/assert"
)

func TestJournalEntryRepository_Create(t *testing.T) {
	// Create a mock DB and expect the named exec
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.JournalEntryRepository{DB: sqlxDB}

	journalEntry := &domain.JournalEntry{
		AccountID:       "1",
		TransactionName: "Test Transaction",
		DebitAmount:     100.0,
		CreditAmount:    100.0,
		Folio:           "Test Folio",
	}

	mock.ExpectExec("INSERT INTO journal_entries \\(account_id, transaction_name, debit_amount, credit_amount, folio\\) VALUES \\(\\?, \\?, \\?, \\?, \\?\\)").
		WithArgs(
			journalEntry.AccountID,
			journalEntry.TransactionName,
			journalEntry.DebitAmount,
			journalEntry.CreditAmount,
			journalEntry.Folio,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Execute the function
	result, err := repo.Create(context.Background(), journalEntry)

	// Assert the expectations
	assert.NoError(t, err)
	assert.Equal(t, journalEntry, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestJournalEntryRepository_Create_Error(t *testing.T) {
	// Create a mock DB and expect the named exec to fail
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.JournalEntryRepository{DB: sqlxDB}

	journalEntry := &domain.JournalEntry{
		AccountID:       "1",
		TransactionName: "Test Transaction",
		DebitAmount:     100.0,
		CreditAmount:    100.0,
		Folio:           "Test Folio",
	}

	mock.ExpectExec("INSERT INTO journal_entries \\(account_id, transaction_name, debit_amount, credit_amount, folio\\) VALUES \\(\\?, \\?, \\?, \\?, \\?\\)").
		WithArgs(
			journalEntry.AccountID,
			journalEntry.TransactionName,
			journalEntry.DebitAmount,
			journalEntry.CreditAmount,
			journalEntry.Folio,
		).
		WillReturnError(errors.New("insert failed"))

	// Execute the function
	result, err := repo.Create(context.Background(), journalEntry)

	// Assert the expectations
	assert.Error(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, &domain.JournalEntry{}, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}
