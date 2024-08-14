package usecase

import (
	"context"
	"database/sql"
	"log"
	"strconv"

	"github.com/krisdioles/ppr-wallet/app/domain"
	"github.com/krisdioles/ppr-wallet/app/domain/errors"
	"github.com/krisdioles/ppr-wallet/app/external"
	"golang.org/x/sync/errgroup"
)

type UserBalanceUsecase struct {
	userBalanceRepository  domain.UserBalanceRepository
	journalEntryRepository domain.JournalEntryRepository
	bank1Client            external.IBank1Client
}

func NewUserBalanceUsecase(userBalanceRepository domain.UserBalanceRepository, journalEntryRepository domain.JournalEntryRepository, bank1Client external.IBank1Client) domain.UserBalanceUsecase {
	return &UserBalanceUsecase{
		userBalanceRepository:  userBalanceRepository,
		journalEntryRepository: journalEntryRepository,
		bank1Client:            bank1Client,
	}
}

func (u *UserBalanceUsecase) GetUserBalanceByID(ctx context.Context, id int64) (*domain.UserBalance, error) {
	userBalance, err := u.userBalanceRepository.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return userBalance, errors.ErrUserNotFound
		}

		return userBalance, err
	}

	return userBalance, nil
}

func (u *UserBalanceUsecase) DisburseBalance(ctx context.Context, id int64) error {
	currentUserBalance, err := u.userBalanceRepository.GetByID(ctx, id)
	if err != nil {
		log.Println("[DisburseBalance] GetByID err:", err)
		if err == sql.ErrNoRows {
			return errors.ErrUserNotFound
		}

		return err
	}

	if currentUserBalance.Balance <= 0 {
		return errors.ErrInsufficientBalance
	}

	// disburse to user's account
	// call external api (bank/3rd party)
	createDisbursementResp, err := u.bank1Client.CreateDisbursement(ctx, &external.Bank1CreateDisbursementRequest{
		ReferenceID: "test-transaction-010121",
		Amount: external.AmountObj{
			Total:    currentUserBalance.Balance,
			Currency: "IDR",
		},
		Account: external.AccountObj{
			AccountHolderName: currentUserBalance.AccountName,
			AccountBankCode:   currentUserBalance.BankCode,
			AccountNo:         currentUserBalance.AccountNo,
		},
	})
	if err != nil {
		log.Println("[DisburseBalance] Create Disbursement err:", err)
		return err
	}

	if createDisbursementResp.Status != "ok" {
		return errors.ErrPartnerError
	}

	if err = u.userBalanceRepository.UpdateBalanceByID(ctx, 0, id); err != nil {
		log.Println("[DisburseBalance] UpdateBalanceByID err:", err)
		return err
	}

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if _, err = u.journalEntryRepository.Create(egCtx, &domain.JournalEntry{
			AccountID:       strconv.Itoa(int(currentUserBalance.ID)),
			TransactionName: "Balance disbursement",
			DebitAmount:     currentUserBalance.Balance,
		}); err != nil {
			log.Println("[DisburseBalance] Create journalentry debit err:", err)
			return err
		}

		if _, err = u.journalEntryRepository.Create(egCtx, &domain.JournalEntry{
			AccountID:       currentUserBalance.AccountNo,
			TransactionName: "Balance disbursement",
			CreditAmount:    currentUserBalance.Balance,
		}); err != nil {
			log.Println("[DisburseBalance] Create journalentry credit err:", err)
			return err
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}
