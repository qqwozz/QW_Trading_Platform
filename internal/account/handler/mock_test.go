package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/qw-trading/platform/internal/models"
)

type mockAccountRepo struct {
	createFn               func(account *models.Account) error
	getByIDFn              func(id uuid.UUID) (*models.Account, error)
	getByUserAndCurrencyFn func(userID uuid.UUID, currency string) (*models.Account, error)
	getByUserIDFn          func(userID uuid.UUID) ([]models.Account, error)
	freezeBalanceFn        func(accountID uuid.UUID, amount float64) error
	unfreezeBalanceFn      func(accountID uuid.UUID, amount float64) error
	creditFn               func(accountID uuid.UUID, amount float64) error
	debitFn                func(accountID uuid.UUID, amount float64) error
	recordBalanceHistoryFn func(entry *models.BalanceHistory) error
}

func (m *mockAccountRepo) Create(_ context.Context, account *models.Account) error {
	if m.createFn != nil {
		return m.createFn(account)
	}
	return nil
}

func (m *mockAccountRepo) GetByID(_ context.Context, id uuid.UUID) (*models.Account, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(id)
	}
	return nil, nil
}

func (m *mockAccountRepo) GetByUserAndCurrency(_ context.Context, userID uuid.UUID, currency string) (*models.Account, error) {
	if m.getByUserAndCurrencyFn != nil {
		return m.getByUserAndCurrencyFn(userID, currency)
	}
	return nil, nil
}

func (m *mockAccountRepo) GetByUserID(_ context.Context, userID uuid.UUID) ([]models.Account, error) {
	if m.getByUserIDFn != nil {
		return m.getByUserIDFn(userID)
	}
	return nil, nil
}

func (m *mockAccountRepo) FreezeBalance(_ context.Context, accountID uuid.UUID, amount float64) error {
	if m.freezeBalanceFn != nil {
		return m.freezeBalanceFn(accountID, amount)
	}
	return nil
}

func (m *mockAccountRepo) UnfreezeBalance(_ context.Context, accountID uuid.UUID, amount float64) error {
	if m.unfreezeBalanceFn != nil {
		return m.unfreezeBalanceFn(accountID, amount)
	}
	return nil
}

func (m *mockAccountRepo) Credit(_ context.Context, accountID uuid.UUID, amount float64) error {
	if m.creditFn != nil {
		return m.creditFn(accountID, amount)
	}
	return nil
}

func (m *mockAccountRepo) Debit(_ context.Context, accountID uuid.UUID, amount float64) error {
	if m.debitFn != nil {
		return m.debitFn(accountID, amount)
	}
	return nil
}

func (m *mockAccountRepo) RecordBalanceHistory(_ context.Context, entry *models.BalanceHistory) error {
	if m.recordBalanceHistoryFn != nil {
		return m.recordBalanceHistoryFn(entry)
	}
	return nil
}
