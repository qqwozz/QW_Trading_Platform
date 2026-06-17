package repository

import (
	"github.com/google/uuid"
	"github.com/qw-trading/platform/internal/models"
)

type OrderRepositoryInterface interface {
	Create(order *models.Order) error
	GetByID(id uuid.UUID) (*models.Order, error)
	UpdateStatus(orderID uuid.UUID, status models.OrderStatus, filledQuantity float64) error
	List(filter ListFilter) ([]models.Order, int, error)
}

type TradeRepositoryInterface interface {
	Create(trade *models.Trade) error
}
