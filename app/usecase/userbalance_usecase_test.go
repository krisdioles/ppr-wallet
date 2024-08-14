package usecase_test

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"testing"

	"github.com/krisdioles/ppr-wallet/app/domain"
	domErr "github.com/krisdioles/ppr-wallet/app/domain/errors"
	"github.com/krisdioles/ppr-wallet/app/domain/mocks"
	"github.com/krisdioles/ppr-wallet/app/external"
	extMocks "github.com/krisdioles/ppr-wallet/app/external/mocks"
	"github.com/krisdioles/ppr-wallet/app/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserBalanceUsecase_GetUserBalanceByID(t *testing.T) {
	ctx := context.Background()
	userID := int64(1)
	expectedUserBalance := &domain.UserBalance{
		ID:      userID,
		Balance: 1000,
	}

	t.Run("Success", func(t *testing.T) {
		mockUserBalanceRepo := new(mocks.UserBalanceRepository)
		usecase := usecase.NewUserBalanceUsecase(mockUserBalanceRepo, nil, nil)

		mockUserBalanceRepo.On("GetByID", ctx, userID).Return(expectedUserBalance, nil)

		result, err := usecase.GetUserBalanceByID(ctx, userID)
		assert.NoError(t, err)
		assert.Equal(t, expectedUserBalance, result)

		mockUserBalanceRepo.AssertExpectations(t)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		mockUserBalanceRepo := new(mocks.UserBalanceRepository)
		usecase := usecase.NewUserBalanceUsecase(mockUserBalanceRepo, nil, nil)

		mockUserBalanceRepo.On("GetByID", ctx, userID).Return(nil, sql.ErrNoRows)

		result, err := usecase.GetUserBalanceByID(ctx, userID)
		assert.ErrorIs(t, err, domErr.ErrUserNotFound)
		assert.Equal(t, (*domain.UserBalance)(nil), result)

		mockUserBalanceRepo.AssertExpectations(t)
	})

	t.Run("OtherError", func(t *testing.T) {
		mockUserBalanceRepo := new(mocks.UserBalanceRepository)
		usecase := usecase.NewUserBalanceUsecase(mockUserBalanceRepo, nil, nil)

		mockUserBalanceRepo.On("GetByID", ctx, userID).Return(nil, errors.New("some error"))

		result, err := usecase.GetUserBalanceByID(ctx, userID)
		assert.Error(t, err)
		assert.Nil(t, result)

		mockUserBalanceRepo.AssertExpectations(t)
	})
}

func TestUserBalanceUsecase_DisburseBalance(t *testing.T) {
	ctx := context.Background()
	userID := int64(1)
	userBalance := &domain.UserBalance{
		ID:          userID,
		Balance:     1000,
		AccountName: "Test User",
		BankCode:    "BANK001",
		AccountNo:   "1234567890",
	}

	t.Run("Success", func(t *testing.T) {
		mockUserBalanceRepo := new(mocks.UserBalanceRepository)
		mockBank1Client := new(extMocks.IBank1Client)
		mockJournalEntryRepo := new(mocks.JournalEntryRepository)

		usecase := usecase.NewUserBalanceUsecase(mockUserBalanceRepo, mockJournalEntryRepo, mockBank1Client)

		mockUserBalanceRepo.On("GetByID", ctx, userID).Return(userBalance, nil)
		mockBank1Client.On("CreateDisbursement", ctx, mock.Anything).Return(&external.Bank1CreateDisbursementResponse{Status: "ok"}, nil)
		mockUserBalanceRepo.On("UpdateBalanceByID", ctx, int64(0), userID).Return(nil)
		mockJournalEntryRepo.On("Create", mock.Anything, &domain.JournalEntry{
			AccountID:       strconv.Itoa(int(userBalance.ID)),
			TransactionName: "Balance disbursement",
			DebitAmount:     userBalance.Balance}).
			Return(&domain.JournalEntry{}, nil)

		mockJournalEntryRepo.On("Create", mock.Anything, &domain.JournalEntry{
			AccountID:       "1234567890",
			TransactionName: "Balance disbursement",
			CreditAmount:    userBalance.Balance}).
			Return(&domain.JournalEntry{}, nil)

		err := usecase.DisburseBalance(ctx, userID)
		assert.NoError(t, err)

		mockUserBalanceRepo.AssertExpectations(t)
		mockBank1Client.AssertExpectations(t)
		mockJournalEntryRepo.AssertExpectations(t)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		mockUserBalanceRepo := new(mocks.UserBalanceRepository)
		mockBank1Client := new(extMocks.IBank1Client)
		mockJournalEntryRepo := new(mocks.JournalEntryRepository)

		usecase := usecase.NewUserBalanceUsecase(mockUserBalanceRepo, mockJournalEntryRepo, mockBank1Client)

		mockUserBalanceRepo.ExpectedCalls = nil

		mockUserBalanceRepo.On("GetByID", ctx, userID).Return(nil, sql.ErrNoRows)

		err := usecase.DisburseBalance(ctx, userID)
		assert.ErrorIs(t, err, domErr.ErrUserNotFound)

		mockUserBalanceRepo.AssertExpectations(t)
	})

	t.Run("InsufficientBalance", func(t *testing.T) {
		mockUserBalanceRepo := new(mocks.UserBalanceRepository)
		mockBank1Client := new(extMocks.IBank1Client)
		mockJournalEntryRepo := new(mocks.JournalEntryRepository)

		usecase := usecase.NewUserBalanceUsecase(mockUserBalanceRepo, mockJournalEntryRepo, mockBank1Client)

		mockUserBalanceRepo.ExpectedCalls = nil

		lowBalance := &domain.UserBalance{
			ID:      userID,
			Balance: 0,
		}

		mockUserBalanceRepo.On("GetByID", ctx, userID).Return(lowBalance, nil)

		err := usecase.DisburseBalance(ctx, userID)
		assert.ErrorIs(t, err, domErr.ErrInsufficientBalance)

		mockUserBalanceRepo.AssertExpectations(t)
	})

	t.Run("CreateDisbursementError", func(t *testing.T) {
		mockUserBalanceRepo := new(mocks.UserBalanceRepository)
		mockBank1Client := new(extMocks.IBank1Client)
		mockJournalEntryRepo := new(mocks.JournalEntryRepository)

		usecase := usecase.NewUserBalanceUsecase(mockUserBalanceRepo, mockJournalEntryRepo, mockBank1Client)

		mockUserBalanceRepo.On("GetByID", ctx, userID).Return(userBalance, nil)
		mockBank1Client.On("CreateDisbursement", ctx, mock.Anything).Return(nil, errors.New("disbursement error"))

		err := usecase.DisburseBalance(ctx, userID)
		assert.Error(t, err)

		mockUserBalanceRepo.AssertExpectations(t)
		mockBank1Client.AssertExpectations(t)
	})

	t.Run("PartnerError", func(t *testing.T) {
		mockUserBalanceRepo := new(mocks.UserBalanceRepository)
		mockBank1Client := new(extMocks.IBank1Client)
		mockJournalEntryRepo := new(mocks.JournalEntryRepository)

		usecase := usecase.NewUserBalanceUsecase(mockUserBalanceRepo, mockJournalEntryRepo, mockBank1Client)

		mockUserBalanceRepo.ExpectedCalls = nil

		mockUserBalanceRepo.On("GetByID", ctx, userID).Return(userBalance, nil)
		mockBank1Client.On("CreateDisbursement", ctx, mock.Anything).Return(&external.Bank1CreateDisbursementResponse{Status: "failed"}, nil)

		err := usecase.DisburseBalance(ctx, userID)
		assert.ErrorIs(t, err, domErr.ErrPartnerError)

		mockUserBalanceRepo.AssertExpectations(t)
		mockBank1Client.AssertExpectations(t)
	})

	t.Run("UpdateBalanceError", func(t *testing.T) {
		mockUserBalanceRepo := new(mocks.UserBalanceRepository)
		mockBank1Client := new(extMocks.IBank1Client)
		mockJournalEntryRepo := new(mocks.JournalEntryRepository)

		usecase := usecase.NewUserBalanceUsecase(mockUserBalanceRepo, mockJournalEntryRepo, mockBank1Client)

		mockUserBalanceRepo.ExpectedCalls = nil

		mockUserBalanceRepo.On("GetByID", ctx, userID).Return(userBalance, nil)
		mockBank1Client.On("CreateDisbursement", ctx, mock.Anything).Return(&external.Bank1CreateDisbursementResponse{Status: "ok"}, nil)
		mockUserBalanceRepo.On("UpdateBalanceByID", ctx, int64(0), userID).Return(errors.New("update error"))

		err := usecase.DisburseBalance(ctx, userID)
		assert.Error(t, err)

		mockUserBalanceRepo.AssertExpectations(t)
		mockBank1Client.AssertExpectations(t)
	})

	t.Run("JournalEntryError", func(t *testing.T) {
		mockUserBalanceRepo := new(mocks.UserBalanceRepository)
		mockBank1Client := new(extMocks.IBank1Client)
		mockJournalEntryRepo := new(mocks.JournalEntryRepository)

		usecase := usecase.NewUserBalanceUsecase(mockUserBalanceRepo, mockJournalEntryRepo, mockBank1Client)

		mockUserBalanceRepo.ExpectedCalls = nil

		mockUserBalanceRepo.On("GetByID", ctx, userID).Return(userBalance, nil)
		mockBank1Client.On("CreateDisbursement", ctx, mock.Anything).Return(&external.Bank1CreateDisbursementResponse{Status: "ok"}, nil)
		mockUserBalanceRepo.On("UpdateBalanceByID", ctx, int64(0), userID).Return(nil)
		mockJournalEntryRepo.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("journal entry error"))
		mockJournalEntryRepo.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("journal entry error"))

		err := usecase.DisburseBalance(ctx, userID)
		assert.Error(t, err)

		mockUserBalanceRepo.AssertExpectations(t)
		mockBank1Client.AssertExpectations(t)
		mockJournalEntryRepo.AssertExpectations(t)
	})
}
