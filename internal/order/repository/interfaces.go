package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/qw-trading/platform/internal/models"
)

type OrderRepositoryInterface interface {
	Create(ctx context.Context, order *models.Order) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Order, error)
	UpdateStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus, filledQuantity float64) error
	List(ctx context.Context, filter ListFilter) ([]models.Order, int, error)
}

type TradeRepositoryInterface interface {
	Create(ctx context.Context, trade *models.Trade) error
}
