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

type Handler struct {
	repo repository.MarketRepositoryInterface
	hub  *hub.Hub
}

func New(repo repository.MarketRepositoryInterface, h *hub.Hub) *Handler {
	return &Handler{repo: repo, hub: h}
}

func (h *Handler) ListTickers(w http.ResponseWriter, r *http.Request) {
	tickers, err := h.repo.GetTickers(r.Context())
	if err != nil {
		response.InternalError(w, "failed to get tickers")
		return
	}

	if tickers == nil {
		tickers = []models.MarketTicker{}
	}

	response.Success(w, tickers)
}

func (h *Handler) GetTicker(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	if symbol == "" {
		response.BadRequest(w, "symbol is required")
		return
	}

	ticker, err := h.repo.GetTicker(r.Context(), symbol)
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

func (h *Handler) GetOrderBook(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	if symbol == "" {
		response.BadRequest(w, "symbol is required")
		return
	}

	depth := 20
	if d, err := strconv.Atoi(r.URL.Query().Get("depth")); err == nil && d > 0 && d <= 1000 {
		depth = d
	}

	snapshot, err := h.repo.GetRecentSnapshot(r.Context(), symbol)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok && appErr.Code == http.StatusNotFound {
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

func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	h.hub.HandleWebSocket(w, r)
}
