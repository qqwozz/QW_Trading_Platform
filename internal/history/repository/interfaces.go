package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/qw-trading/platform/internal/models"
)

type HistoryRepositoryInterface interface {
	RecordBalanceChange(ctx context.Context, entry *models.BalanceHistory) error
	RecordPositionChange(ctx context.Context, entry *models.PositionHistory) error
	GetOrderHistory(ctx context.Context, userID uuid.UUID, symbol string, limit, offset int) ([]models.Order, int, error)
	GetTradeHistory(ctx context.Context, userID uuid.UUID, symbol string, limit, offset int) ([]models.Trade, int, error)
	GetBalanceHistory(ctx context.Context, userID uuid.UUID, currency string, limit, offset int) ([]models.BalanceHistory, int, error)
	GetPositionHistory(ctx context.Context, userID uuid.UUID, symbol string, limit, offset int) ([]models.PositionHistory, int, error)
}
