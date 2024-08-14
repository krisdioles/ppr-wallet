package provider

import (
	"github.com/jmoiron/sqlx"
	"github.com/krisdioles/ppr-wallet/app/domain"
	"github.com/krisdioles/ppr-wallet/app/repository"
)

type Repository struct {
	UserBalanceRepository  domain.UserBalanceRepository
	JournalEntryRepository domain.JournalEntryRepository
}

func InitRepositories(db *sqlx.DB) *Repository {
	return &Repository{
		UserBalanceRepository:  repository.NewUserBalanceRepository(db),
		JournalEntryRepository: repository.NewJournalEntryRepository(db),
	}
}
