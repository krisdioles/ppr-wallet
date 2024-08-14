package domain

import "context"

type JournalEntry struct {
	ID              int64  `json:"id" db:"id"`
	AccountID       string `json:"account_id" db:"account_id"`
	TransactionName string `json:"transaction_name" db:"transaction_name"`
	DebitAmount     int64  `json:"debit_amount" db:"debit_amount"`
	CreditAmount    int64  `json:"credit_amount" db:"credit_amount"`
	Folio           string `json:"folio" db:"folio"`
}

func (j *JournalEntry) TableName() string {
	return "journal_entries"
}

type JournalEntryRepository interface {
	Create(ctx context.Context, journalEntry *JournalEntry) (*JournalEntry, error)
}
