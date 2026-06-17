package handler

import (
	"github.com/google/uuid"
	"github.com/qw-trading/platform/internal/models"
	"github.com/qw-trading/platform/internal/order/repository"
)

type mockOrderRepo struct {
	createFn      func(order *models.Order) error
	getByIDFn     func(id uuid.UUID) (*models.Order, error)
	updateStatusFn func(orderID uuid.UUID, status models.OrderStatus, filledQuantity float64) error
	listFn        func(filter repository.ListFilter) ([]models.Order, int, error)
}

func (m *mockOrderRepo) Create(order *models.Order) error {
	if m.createFn != nil {
		return m.createFn(order)
	}
	return nil
}

func (m *mockOrderRepo) GetByID(id uuid.UUID) (*models.Order, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(id)
	}
	return nil, nil
}

func (m *mockOrderRepo) UpdateStatus(orderID uuid.UUID, status models.OrderStatus, filledQuantity float64) error {
	if m.updateStatusFn != nil {
		return m.updateStatusFn(orderID, status, filledQuantity)
	}
	return nil
}

func (m *mockOrderRepo) List(filter repository.ListFilter) ([]models.Order, int, error) {
	if m.listFn != nil {
		return m.listFn(filter)
	}
	return nil, 0, nil
}

type mockTradeRepo struct {
	createFn func(trade *models.Trade) error
}

func (m *mockTradeRepo) Create(trade *models.Trade) error {
	if m.createFn != nil {
		return m.createFn(trade)
	}
	return nil
}
