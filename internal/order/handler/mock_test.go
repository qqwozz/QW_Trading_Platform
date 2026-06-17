package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/qw-trading/platform/internal/models"
	"github.com/qw-trading/platform/internal/order/repository"
)

type mockOrderRepo struct {
	createFn       func(order *models.Order) error
	getByIDFn      func(id uuid.UUID) (*models.Order, error)
	updateStatusFn func(orderID uuid.UUID, status models.OrderStatus, filledQuantity float64) error
	listFn         func(filter repository.ListFilter) ([]models.Order, int, error)
}

func (m *mockOrderRepo) Create(_ context.Context, order *models.Order) error {
	if m.createFn != nil {
		return m.createFn(order)
	}
	return nil
}

func (m *mockOrderRepo) GetByID(_ context.Context, id uuid.UUID) (*models.Order, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(id)
	}
	return nil, nil
}

func (m *mockOrderRepo) UpdateStatus(_ context.Context, orderID uuid.UUID, status models.OrderStatus, filledQuantity float64) error {
	if m.updateStatusFn != nil {
		return m.updateStatusFn(orderID, status, filledQuantity)
	}
	return nil
}

func (m *mockOrderRepo) List(_ context.Context, filter repository.ListFilter) ([]models.Order, int, error) {
	if m.listFn != nil {
		return m.listFn(filter)
	}
	return nil, 0, nil
}

type mockTradeRepo struct {
	createFn func(trade *models.Trade) error
}

func (m *mockTradeRepo) Create(_ context.Context, trade *models.Trade) error {
	if m.createFn != nil {
		return m.createFn(trade)
	}
	return nil
}
