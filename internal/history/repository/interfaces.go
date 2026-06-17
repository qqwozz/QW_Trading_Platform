package repository

import (
	"github.com/google/uuid"
	"github.com/qw-trading/platform/internal/models"
)

type HistoryRepositoryInterface interface {
	RecordBalanceChange(entry *models.BalanceHistory) error
	RecordPositionChange(entry *models.PositionHistory) error
	GetOrderHistory(userID uuid.UUID, symbol string, limit, offset int) ([]models.Order, int, error)
	GetTradeHistory(userID uuid.UUID, symbol string, limit, offset int) ([]models.Trade, int, error)
	GetBalanceHistory(userID uuid.UUID, currency string, limit, offset int) ([]models.BalanceHistory, int, error)
	GetPositionHistory(userID uuid.UUID, symbol string, limit, offset int) ([]models.PositionHistory, int, error)
}
