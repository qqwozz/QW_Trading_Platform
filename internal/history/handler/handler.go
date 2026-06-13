package handler

import (
	"net/http"
	"strconv"

	"github.com/qw-trading/platform/internal/history/repository"
	"github.com/qw-trading/platform/internal/models"
	"github.com/qw-trading/platform/pkg/middleware"
	"github.com/qw-trading/platform/pkg/response"
)

type Handler struct {
	repo *repository.HistoryRepository
}

func New(repo *repository.HistoryRepository) *Handler {
	return &Handler{repo: repo}
}

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

func (h *Handler) GetOrderHistory(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Unauthorized(w, "unauthorized")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	orders, total, err := h.repo.GetOrderHistory(userID, r.URL.Query().Get("symbol"), limit, offset)
	if err != nil {
		response.InternalError(w, "failed to get order history")
		return
	}

	var result []OrderResponse
	for _, order := range orders {
		result = append(result, toOrderResponse(&order))
	}

	response.Paginated(w, result, total, limit, offset)
}

func (h *Handler) GetTradeHistory(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Unauthorized(w, "unauthorized")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	trades, total, err := h.repo.GetTradeHistory(userID, r.URL.Query().Get("symbol"), limit, offset)
	if err != nil {
		response.InternalError(w, "failed to get trade history")
		return
	}

	var result []TradeResponse
	for _, trade := range trades {
		result = append(result, toTradeResponse(&trade))
	}

	response.Paginated(w, result, total, limit, offset)
}

func (h *Handler) GetBalanceHistory(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Unauthorized(w, "unauthorized")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	history, total, err := h.repo.GetBalanceHistory(userID, r.URL.Query().Get("currency"), limit, offset)
	if err != nil {
		response.InternalError(w, "failed to get balance history")
		return
	}

	var result []BalanceHistoryResponse
	for _, entry := range history {
		result = append(result, toBalanceHistoryResponse(&entry))
	}

	response.Paginated(w, result, total, limit, offset)
}

func (h *Handler) GetPositionHistory(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Unauthorized(w, "unauthorized")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	history, total, err := h.repo.GetPositionHistory(userID, r.URL.Query().Get("symbol"), limit, offset)
	if err != nil {
		response.InternalError(w, "failed to get position history")
		return
	}

	var result []PositionHistoryResponse
	for _, entry := range history {
		result = append(result, toPositionHistoryResponse(&entry))
	}

	response.Paginated(w, result, total, limit, offset)
}

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
