package repository

import (
	"github.com/google/uuid"
	"github.com/qw-trading/platform/internal/models"
)

type AccountRepositoryInterface interface {
	Create(account *models.Account) error
	GetByID(id uuid.UUID) (*models.Account, error)
	GetByUserAndCurrency(userID uuid.UUID, currency string) (*models.Account, error)
	GetByUserID(userID uuid.UUID) ([]models.Account, error)
	FreezeBalance(accountID uuid.UUID, amount float64) error
	UnfreezeBalance(accountID uuid.UUID, amount float64) error
	Credit(accountID uuid.UUID, amount float64) error
	Debit(accountID uuid.UUID, amount float64) error
	RecordBalanceHistory(entry *models.BalanceHistory) error
}
