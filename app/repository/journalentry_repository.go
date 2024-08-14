package repository

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/krisdioles/ppr-wallet/app/domain"
)

type JournalEntryRepository struct {
	DB *sqlx.DB
}

func NewJournalEntryRepository(db *sqlx.DB) *JournalEntryRepository {
	return &JournalEntryRepository{
		DB: db,
	}
}

func (r *JournalEntryRepository) Create(ctx context.Context, journalEntry *domain.JournalEntry) (*domain.JournalEntry, error) {
	createJournalEntryQuery := `INSERT INTO journal_entries 
	(account_id, transaction_name, debit_amount, credit_amount, folio) VALUES
	(:account_id, :transaction_name, :debit_amount, :credit_amount, :folio)`

	result, err := r.DB.NamedExecContext(ctx, createJournalEntryQuery, journalEntry)
	if err != nil {
		return &domain.JournalEntry{}, err
	}
	log.Println("[Create] result:", result)

	return journalEntry, nil
}
