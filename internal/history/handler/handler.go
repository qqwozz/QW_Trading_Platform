// Package handler implements HTTP handlers for history and audit trail endpoints.
package handler

import (
	"net/http"

	"github.com/qw-trading/platform/internal/history/repository"
	"github.com/qw-trading/platform/internal/models"
	"github.com/qw-trading/platform/pkg/middleware"
	"github.com/qw-trading/platform/pkg/response"
)

// Handler holds dependencies for history-related HTTP handlers.
type Handler struct {
	repo repository.HistoryRepositoryInterface
}

// New creates a new Handler with the given repository.
func New(repo repository.HistoryRepositoryInterface) *Handler {
	return &Handler{repo: repo}
}

// OrderResponse is the JSON response containing order history information.
type OrderResponse struct {
	ID             string   `json:"id"`
	Symbol         string   `json:"symbol"`
	Side           string   `json:"side"`
	Type           string   `json:"type"`
	Price          *float64 `json:"price,omitempty"`
	Quantity       float64  `json:"quantity"`
	FilledQuantity float64  `json:"filled_quantity"`
	Status         string   `json:"status"`
	TimeInForce    string   `json:"time_in_force"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
}

// TradeResponse is the JSON response containing trade history information.
type TradeResponse struct {
	ID            string  `json:"id"`
	Symbol        string  `json:"symbol"`
	BuyerOrderID  string  `json:"buyer_order_id"`
	SellerOrderID string  `json:"seller_order_id"`
	BuyerID       string  `json:"buyer_id"`
	SellerID      string  `json:"seller_id"`
	Price         float64 `json:"price"`
	Quantity      float64 `json:"quantity"`
	BuyerFee      float64 `json:"buyer_fee"`
	SellerFee     float64 `json:"seller_fee"`
	ExecutedAt    string  `json:"executed_at"`
}

// BalanceHistoryResponse is the JSON response containing balance change history.
type BalanceHistoryResponse struct {
	ID            string  `json:"id"`
	Currency      string  `json:"currency"`
	Amount        float64 `json:"amount"`
	BalanceBefore float64 `json:"balance_before"`
	BalanceAfter  float64 `json:"balance_after"`
	Type          string  `json:"type"`
	Description   string  `json:"description"`
	CreatedAt     string  `json:"created_at"`
}

// PositionHistoryResponse is the JSON response containing position change history.
type PositionHistoryResponse struct {
	ID             string  `json:"id"`
	Symbol         string  `json:"symbol"`
	QuantityChange float64 `json:"quantity_change"`
	QuantityBefore float64 `json:"quantity_before"`
	QuantityAfter  float64 `json:"quantity_after"`
	AvgPriceBefore float64 `json:"avg_price_before"`
	AvgPriceAfter  float64 `json:"avg_price_after"`
	Type           string  `json:"type"`
	TradeID        string  `json:"trade_id,omitempty"`
	CreatedAt      string  `json:"created_at"`
}

// GetOrderHistory handles GET /history/orders. Returns paginated order history
// for the authenticated user, optionally filtered by symbol.
func (h *Handler) GetOrderHistory(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Unauthorized(w, "unauthorized")
		return
	}

	limit, offset := response.ParsePagination(r, 50)

	orders, total, err := h.repo.GetOrderHistory(userID, r.URL.Query().Get("symbol"), limit, offset)
	if err != nil {
		response.InternalError(w, "failed to get order history")
		return
	}

	result := make([]OrderResponse, len(orders))
	for i, order := range orders {
		result[i] = toOrderResponse(&order)
	}

	response.Paginated(w, result, total, limit, offset)
}

// GetTradeHistory handles GET /history/trades. Returns paginated trade history
// for the authenticated user, optionally filtered by symbol.
func (h *Handler) GetTradeHistory(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Unauthorized(w, "unauthorized")
		return
	}

	limit, offset := response.ParsePagination(r, 50)

	trades, total, err := h.repo.GetTradeHistory(userID, r.URL.Query().Get("symbol"), limit, offset)
	if err != nil {
		response.InternalError(w, "failed to get trade history")
		return
	}

	result := make([]TradeResponse, len(trades))
	for i, trade := range trades {
		result[i] = toTradeResponse(&trade)
	}

	response.Paginated(w, result, total, limit, offset)
}

// GetBalanceHistory handles GET /history/balance. Returns paginated balance
// change history for the authenticated user, optionally filtered by currency.
func (h *Handler) GetBalanceHistory(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Unauthorized(w, "unauthorized")
		return
	}

	limit, offset := response.ParsePagination(r, 50)

	history, total, err := h.repo.GetBalanceHistory(userID, r.URL.Query().Get("currency"), limit, offset)
	if err != nil {
		response.InternalError(w, "failed to get balance history")
		return
	}

	result := make([]BalanceHistoryResponse, len(history))
	for i, entry := range history {
		result[i] = toBalanceHistoryResponse(&entry)
	}

	response.Paginated(w, result, total, limit, offset)
}

// GetPositionHistory handles GET /history/positions. Returns paginated position
// change history for the authenticated user, optionally filtered by symbol.
func (h *Handler) GetPositionHistory(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Unauthorized(w, "unauthorized")
		return
	}

	limit, offset := response.ParsePagination(r, 50)

	history, total, err := h.repo.GetPositionHistory(userID, r.URL.Query().Get("symbol"), limit, offset)
	if err != nil {
		response.InternalError(w, "failed to get position history")
		return
	}

	result := make([]PositionHistoryResponse, len(history))
	for i, entry := range history {
		result[i] = toPositionHistoryResponse(&entry)
	}

	response.Paginated(w, result, total, limit, offset)
}

// toOrderResponse converts a domain Order into an API OrderResponse.
func toOrderResponse(order *models.Order) OrderResponse {
	return OrderResponse{
		ID:             order.ID.String(),
		Symbol:         order.Symbol,
		Side:           string(order.Side),
		Type:           string(order.Type),
		Price:          order.Price,
		Quantity:       order.Quantity,
		FilledQuantity: order.FilledQuantity,
		Status:         string(order.Status),
		TimeInForce:    string(order.TimeInForce),
		CreatedAt:      order.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:      order.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// toTradeResponse converts a domain Trade into an API TradeResponse.
func toTradeResponse(trade *models.Trade) TradeResponse {
	return TradeResponse{
		ID:            trade.ID.String(),
		Symbol:        trade.Symbol,
		BuyerOrderID:  trade.BuyerOrderID.String(),
		SellerOrderID: trade.SellerOrderID.String(),
		BuyerID:       trade.BuyerID.String(),
		SellerID:      trade.SellerID.String(),
		Price:         trade.Price,
		Quantity:      trade.Quantity,
		BuyerFee:      trade.BuyerFee,
		SellerFee:     trade.SellerFee,
		ExecutedAt:    trade.ExecutedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// toBalanceHistoryResponse converts a domain BalanceHistory into an API BalanceHistoryResponse.
func toBalanceHistoryResponse(entry *models.BalanceHistory) BalanceHistoryResponse {
	return BalanceHistoryResponse{
		ID:            entry.ID.String(),
		Currency:      entry.Currency,
		Amount:        entry.Amount,
		BalanceBefore: entry.BalanceBefore,
		BalanceAfter:  entry.BalanceAfter,
		Type:          string(entry.Type),
		Description:   entry.Description,
		CreatedAt:     entry.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// toPositionHistoryResponse converts a domain PositionHistory into an API PositionHistoryResponse.
func toPositionHistoryResponse(entry *models.PositionHistory) PositionHistoryResponse {
	resp := PositionHistoryResponse{
		ID:             entry.ID.String(),
		Symbol:         entry.Symbol,
		QuantityChange: entry.QuantityChange,
		QuantityBefore: entry.QuantityBefore,
		QuantityAfter:  entry.QuantityAfter,
		AvgPriceBefore: entry.AvgPriceBefore,
		AvgPriceAfter:  entry.AvgPriceAfter,
		Type:           string(entry.Type),
		CreatedAt:      entry.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if entry.TradeID != nil {
		resp.TradeID = entry.TradeID.String()
	}
	return resp
}
