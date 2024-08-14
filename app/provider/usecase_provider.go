package provider

import (
	"github.com/krisdioles/ppr-wallet/app/domain"
	"github.com/krisdioles/ppr-wallet/app/external"
	"github.com/krisdioles/ppr-wallet/app/usecase"
	"github.com/krisdioles/ppr-wallet/config"
)

type Usecase struct {
	UserBalanceUsecase domain.UserBalanceUsecase
	Bank1Client        external.Bank1Client
}

func InitUsecases(cfg *config.Config, repo *Repository) *Usecase {
	bank1Client := external.NewBank1Client(&cfg.Bank1)

	return &Usecase{
		UserBalanceUsecase: usecase.NewUserBalanceUsecase(repo.UserBalanceRepository, repo.JournalEntryRepository, bank1Client),
	}
}
