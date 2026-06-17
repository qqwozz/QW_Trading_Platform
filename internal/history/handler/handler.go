package handler

import (
	"net/http"

	"github.com/qw-trading/platform/internal/history/repository"
	"github.com/qw-trading/platform/internal/models"
	"github.com/qw-trading/platform/pkg/middleware"
	"github.com/qw-trading/platform/pkg/response"
)

type Handler struct {
	repo repository.HistoryRepositoryInterface
}

func New(repo repository.HistoryRepositoryInterface) *Handler {
	return &Handler{repo: repo}
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

	limit, offset := response.ParsePagination(r, 50)

	orders, total, err := h.repo.GetOrderHistory(r.Context(), userID, r.URL.Query().Get("symbol"), limit, offset)
	if err != nil {
		response.InternalError(w, "failed to get order history")
		return
	}

	result := make([]response.OrderResponse, len(orders))
	for i, order := range orders {
		result[i] = response.OrderFromModel(&order)
	}

	response.Paginated(w, result, total, limit, offset)
}

func (h *Handler) GetTradeHistory(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Unauthorized(w, "unauthorized")
		return
	}

	limit, offset := response.ParsePagination(r, 50)

	trades, total, err := h.repo.GetTradeHistory(r.Context(), userID, r.URL.Query().Get("symbol"), limit, offset)
	if err != nil {
		response.InternalError(w, "failed to get trade history")
		return
	}

	result := make([]response.TradeResponse, len(trades))
	for i, trade := range trades {
		result[i] = response.TradeFromModel(&trade)
	}

	response.Paginated(w, result, total, limit, offset)
}

func (h *Handler) GetBalanceHistory(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Unauthorized(w, "unauthorized")
		return
	}

	limit, offset := response.ParsePagination(r, 50)

	history, total, err := h.repo.GetBalanceHistory(r.Context(), userID, r.URL.Query().Get("currency"), limit, offset)
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

func (h *Handler) GetPositionHistory(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Unauthorized(w, "unauthorized")
		return
	}

	limit, offset := response.ParsePagination(r, 50)

	history, total, err := h.repo.GetPositionHistory(r.Context(), userID, r.URL.Query().Get("symbol"), limit, offset)
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
