// Package handler implements HTTP handlers for market data endpoints including
// tickers, order books, and WebSocket connections.
package handler

import (
	"net/http"
	"strconv"

	"github.com/qw-trading/platform/internal/market/hub"
	"github.com/qw-trading/platform/internal/market/repository"
	"github.com/qw-trading/platform/internal/models"
	"github.com/qw-trading/platform/pkg/errors"
	"github.com/qw-trading/platform/pkg/response"
)

// Handler holds dependencies for market data HTTP handlers.
type Handler struct {
	repo *repository.MarketRepository
	hub  *hub.Hub
}

// New creates a new Handler with the given repository and WebSocket hub.
func New(repo *repository.MarketRepository, h *hub.Hub) *Handler {
	return &Handler{repo: repo, hub: h}
}

// ListTickers handles GET /market/tickers. Returns all available market tickers.
func (h *Handler) ListTickers(w http.ResponseWriter, r *http.Request) {
	tickers, err := h.repo.GetTickers()
	if err != nil {
		response.InternalError(w, "failed to get tickers")
		return
	}

	// Ensure a non-null slice is returned in the JSON response.
	if tickers == nil {
		tickers = []models.MarketTicker{}
	}

	response.Success(w, tickers)
}

// GetTicker handles GET /market/tickers/{symbol}. Returns the ticker for
// the specified trading pair.
func (h *Handler) GetTicker(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	if symbol == "" {
		response.BadRequest(w, "symbol is required")
		return
	}

	ticker, err := h.repo.GetTicker(symbol)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok && appErr.Code == http.StatusNotFound {
			response.NotFound(w, "ticker not found")
			return
		}
		response.InternalError(w, "failed to get ticker")
		return
	}

	response.Success(w, ticker)
}

// GetOrderBook handles GET /market/orderbook/{symbol}. Returns the order book
// snapshot limited to the requested depth (default 20).
func (h *Handler) GetOrderBook(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	if symbol == "" {
		response.BadRequest(w, "symbol is required")
		return
	}

	depthStr := r.URL.Query().Get("depth")
	depth := 20
	if depthStr != "" {
		if d, err := strconv.Atoi(depthStr); err == nil && d > 0 {
			depth = d
		}
	}

	snapshot, err := h.repo.GetRecentSnapshot(symbol)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok && appErr.Code == http.StatusNotFound {
			// Return empty order book rather than an error when no snapshot exists.
			response.Success(w, models.OrderBookResponse{
				Symbol: symbol,
				Bids:   []models.OrderBookLevel{},
				Asks:   []models.OrderBookLevel{},
				Depth:  depth,
			})
			return
		}
		response.InternalError(w, "failed to get order book")
		return
	}

	// Trim bids and asks to the requested depth.
	bids := make([]models.OrderBookLevel, 0, min(depth, len(snapshot.Bids)))
	for i, b := range snapshot.Bids {
		if i >= depth {
			break
		}
		bids = append(bids, models.OrderBookLevel{Price: b.Price, Quantity: b.Quantity})
	}

	asks := make([]models.OrderBookLevel, 0, min(depth, len(snapshot.Asks)))
	for i, a := range snapshot.Asks {
		if i >= depth {
			break
		}
		asks = append(asks, models.OrderBookLevel{Price: a.Price, Quantity: a.Quantity})
	}

	response.Success(w, models.OrderBookResponse{
		Symbol: symbol,
		Bids:   bids,
		Asks:   asks,
		Depth:  depth,
	})
}

// HandleWebSocket handles GET /market/ws. Upgrades the connection to WebSocket
// and registers it with the hub for real-time market data.
func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	h.hub.HandleWebSocket(w, r)
}

// min returns the smaller of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
