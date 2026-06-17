package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/qw-trading/platform/internal/models"
	"github.com/qw-trading/platform/internal/order/repository"
	"github.com/qw-trading/platform/pkg/middleware"
)

func TestCreateOrder_Success(t *testing.T) {
	userID := uuid.New()
	orderRepo := &mockOrderRepo{
		createFn: func(order *models.Order) error {
			order.CreatedAt = time.Now()
			order.UpdatedAt = time.Now()
			return nil
		},
	}
	tradeRepo := &mockTradeRepo{}
	h := New(orderRepo, tradeRepo)

	price := 50000.0
	body, _ := json.Marshal(CreateOrderRequest{
		Symbol:   "BTC/USDT",
		Side:     "BUY",
		Type:     "LIMIT",
		Price:    &price,
		Quantity: 0.1,
	})
	req := httptest.NewRequest("POST", "/orders", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID.String())
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	h.CreateOrder(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", w.Code, http.StatusCreated)
	}
}

func TestCreateOrder_InvalidSide(t *testing.T) {
	repo := &mockOrderRepo{}
	tradeRepo := &mockTradeRepo{}
	h := New(repo, tradeRepo)

	body, _ := json.Marshal(CreateOrderRequest{
		Symbol:   "BTC/USDT",
		Side:     "INVALID",
		Type:     "LIMIT",
		Quantity: 0.1,
	})
	req := httptest.NewRequest("POST", "/orders", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, uuid.New().String())
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	h.CreateOrder(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCreateOrder_MarketWithoutPrice(t *testing.T) {
	repo := &mockOrderRepo{}
	tradeRepo := &mockTradeRepo{}
	h := New(repo, tradeRepo)

	body, _ := json.Marshal(CreateOrderRequest{
		Symbol:   "BTC/USDT",
		Side:     "BUY",
		Type:     "MARKET",
		Quantity: 0.1,
	})
	req := httptest.NewRequest("POST", "/orders", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, uuid.New().String())
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	h.CreateOrder(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d (market orders don't need price)", w.Code, http.StatusCreated)
	}
}

func TestCreateOrder_LimitWithoutPrice(t *testing.T) {
	repo := &mockOrderRepo{}
	tradeRepo := &mockTradeRepo{}
	h := New(repo, tradeRepo)

	body, _ := json.Marshal(CreateOrderRequest{
		Symbol:   "BTC/USDT",
		Side:     "BUY",
		Type:     "LIMIT",
		Quantity: 0.1,
	})
	req := httptest.NewRequest("POST", "/orders", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, uuid.New().String())
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	h.CreateOrder(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d (limit requires price)", w.Code, http.StatusBadRequest)
	}
}

func TestCreateOrder_NegativeQuantity(t *testing.T) {
	repo := &mockOrderRepo{}
	tradeRepo := &mockTradeRepo{}
	h := New(repo, tradeRepo)

	body, _ := json.Marshal(CreateOrderRequest{
		Symbol:   "BTC/USDT",
		Side:     "BUY",
		Type:     "MARKET",
		Quantity: -1,
	})
	req := httptest.NewRequest("POST", "/orders", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, uuid.New().String())
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	h.CreateOrder(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGetOrder_Success(t *testing.T) {
	userID := uuid.New()
	orderID := uuid.New()
	orderRepo := &mockOrderRepo{
		getByIDFn: func(id uuid.UUID) (*models.Order, error) {
			return &models.Order{
				ID:        id,
				UserID:    userID,
				Symbol:    "BTC/USDT",
				Side:      models.OrderSideBuy,
				Type:      models.OrderTypeLimit,
				Quantity:  1,
				Status:    models.OrderStatusOpen,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		},
	}
	h := New(orderRepo, &mockTradeRepo{})

	req := httptest.NewRequest("GET", "/orders/"+orderID.String(), nil)
	req.SetPathValue("id", orderID.String())
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID.String())
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	h.GetOrder(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetOrder_Forbidden(t *testing.T) {
	userID := uuid.New()
	orderID := uuid.New()
	otherUserID := uuid.New()
	orderRepo := &mockOrderRepo{
		getByIDFn: func(id uuid.UUID) (*models.Order, error) {
			return &models.Order{
				ID:     id,
				UserID: otherUserID,
				Status: models.OrderStatusOpen,
			}, nil
		},
	}
	h := New(orderRepo, &mockTradeRepo{})

	req := httptest.NewRequest("GET", "/orders/"+orderID.String(), nil)
	req.SetPathValue("id", orderID.String())
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID.String())
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	h.GetOrder(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("status = %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestListOrders_Success(t *testing.T) {
	userID := uuid.New()
	orderRepo := &mockOrderRepo{
		listFn: func(filter repository.ListFilter) ([]models.Order, int, error) {
			return []models.Order{
				{ID: uuid.New(), UserID: userID, Symbol: "BTC/USDT", Status: models.OrderStatusOpen, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			}, 1, nil
		},
	}
	h := New(orderRepo, &mockTradeRepo{})

	req := httptest.NewRequest("GET", "/orders", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID.String())
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	h.ListOrders(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestCancelOrder_Success(t *testing.T) {
	userID := uuid.New()
	orderID := uuid.New()
	orderRepo := &mockOrderRepo{
		getByIDFn: func(id uuid.UUID) (*models.Order, error) {
			return &models.Order{
				ID:        id,
				UserID:    userID,
				Status:    models.OrderStatusOpen,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		},
		updateStatusFn: func(orderID uuid.UUID, status models.OrderStatus, filledQuantity float64) error {
			return nil
		},
	}
	h := New(orderRepo, &mockTradeRepo{})

	req := httptest.NewRequest("DELETE", "/orders/"+orderID.String(), nil)
	req.SetPathValue("id", orderID.String())
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID.String())
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	h.CancelOrder(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestCancelOrder_AlreadyFilled(t *testing.T) {
	userID := uuid.New()
	orderID := uuid.New()
	orderRepo := &mockOrderRepo{
		getByIDFn: func(id uuid.UUID) (*models.Order, error) {
			return &models.Order{
				ID:     id,
				UserID: userID,
				Status: models.OrderStatusFilled,
			}, nil
		},
	}
	h := New(orderRepo, &mockTradeRepo{})

	req := httptest.NewRequest("DELETE", "/orders/"+orderID.String(), nil)
	req.SetPathValue("id", orderID.String())
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID.String())
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	h.CancelOrder(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("status = %d, want %d", w.Code, http.StatusConflict)
	}
}
