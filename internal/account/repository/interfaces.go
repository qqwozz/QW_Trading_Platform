package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/qw-trading/platform/internal/models"
)

type AccountRepositoryInterface interface {
	Create(ctx context.Context, account *models.Account) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Account, error)
	GetByUserAndCurrency(ctx context.Context, userID uuid.UUID, currency string) (*models.Account, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.Account, error)
	FreezeBalance(ctx context.Context, accountID uuid.UUID, amount float64) error
	UnfreezeBalance(ctx context.Context, accountID uuid.UUID, amount float64) error
	Credit(ctx context.Context, accountID uuid.UUID, amount float64) error
	Debit(ctx context.Context, accountID uuid.UUID, amount float64) error
	RecordBalanceHistory(ctx context.Context, entry *models.BalanceHistory) error
}
