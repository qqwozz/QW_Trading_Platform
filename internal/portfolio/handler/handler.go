package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/qw-trading/platform/internal/models"
	"github.com/qw-trading/platform/internal/portfolio/repository"
	"github.com/qw-trading/platform/pkg/errors"
	"github.com/qw-trading/platform/pkg/middleware"
	"github.com/qw-trading/platform/pkg/response"
)

type Handler struct {
	posRepo *repository.PositionRepository
	accRepo *repository.AccountRepository
}

func New(posRepo *repository.PositionRepository, accRepo *repository.AccountRepository) *Handler {
	return &Handler{posRepo: posRepo, accRepo: accRepo}
}

type PositionResponse struct {
	ID            string  `json:"id"`
	Symbol        string  `json:"symbol"`
	Quantity      float64 `json:"quantity"`
	AveragePrice  float64 `json:"average_price"`
	CurrentPrice  float64 `json:"current_price"`
	UnrealizedPnL float64 `json:"unrealized_pnl"`
	CostBasis     float64 `json:"cost_basis"`
	MarketValue   float64 `json:"market_value"`
	UpdatedAt     string  `json:"updated_at"`
}

type PortfolioSummaryResponse struct {
	TotalBalance       float64            `json:"total_balance"`
	TotalFrozen        float64            `json:"total_frozen"`
	TotalMarketValue   float64            `json:"total_market_value"`
	TotalUnrealizedPnL float64            `json:"total_unrealized_pnl"`
	TotalEquity        float64            `json:"total_equity"`
	Currency           string             `json:"currency"`
	Positions          []PositionResponse `json:"positions"`
}

type BalanceResponse struct {
	Currency  string  `json:"currency"`
	Balance   float64 `json:"balance"`
	Frozen    float64 `json:"frozen"`
	Available float64 `json:"available"`
}

func (h *Handler) GetPortfolio(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Unauthorized(w, "unauthorized")
		return
	}

	accounts, err := h.accRepo.GetByUserID(userID)
	if err != nil {
		response.InternalError(w, "failed to get accounts")
		return
	}

	var totalBalance, totalFrozen float64
	for _, acc := range accounts {
		if acc.Currency == "USDT" {
			totalBalance += acc.Balance
			totalFrozen += acc.FrozenBalance
		}
	}

	positions, err := h.posRepo.GetByUserID(userID)
	if err != nil {
		response.InternalError(w, "failed to get positions")
		return
	}

	var totalMarketValue, totalUnrealizedPnL float64
	var posResponses []PositionResponse

	for _, pos := range positions {
		resp := toPositionResponse(&pos)
		totalMarketValue += resp.MarketValue
		totalUnrealizedPnL += resp.UnrealizedPnL
		posResponses = append(posResponses, resp)
	}

	response.Success(w, PortfolioSummaryResponse{
		TotalBalance:       totalBalance,
		TotalFrozen:        totalFrozen,
		TotalMarketValue:   totalMarketValue,
		TotalUnrealizedPnL: totalUnrealizedPnL,
		TotalEquity:        totalBalance + totalMarketValue + totalUnrealizedPnL,
		Currency:           "USDT",
		Positions:          posResponses,
	})
}

func (h *Handler) ListPositions(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Unauthorized(w, "unauthorized")
		return
	}

	positions, err := h.posRepo.GetByUserID(userID)
	if err != nil {
		response.InternalError(w, "failed to get positions")
		return
	}

	var result []PositionResponse
	for _, pos := range positions {
		result = append(result, toPositionResponse(&pos))
	}

	response.Success(w, map[string]interface{}{"positions": result})
}

func (h *Handler) GetBalances(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Unauthorized(w, "unauthorized")
		return
	}

	accounts, err := h.accRepo.GetByUserID(userID)
	if err != nil {
		response.InternalError(w, "failed to get accounts")
		return
	}

	var result []BalanceResponse
	for _, acc := range accounts {
		result = append(result, BalanceResponse{
			Currency:  acc.Currency,
			Balance:   acc.Balance,
			Frozen:    acc.FrozenBalance,
			Available: acc.Balance - acc.FrozenBalance,
		})
	}

	response.Success(w, map[string]interface{}{"balances": result})
}

type UpdatePositionRequest struct {
	Symbol   string  `json:"symbol"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
	Side     string  `json:"side"`
}

func (h *Handler) UpdatePosition(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Unauthorized(w, "unauthorized")
		return
	}

	var req UpdatePositionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if req.Symbol == "" || req.Quantity <= 0 || req.Price <= 0 {
		response.BadRequest(w, "symbol, positive quantity, and positive price are required")
		return
	}

	if req.Side != "BUY" && req.Side != "SELL" {
		response.BadRequest(w, "side must be BUY or SELL")
		return
	}

	pos, err := h.posRepo.GetByUserAndSymbol(userID, req.Symbol)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok && appErr.Code == http.StatusNotFound {
			pos = &models.Position{
				ID:        uuid.New(),
				UserID:    userID,
				Symbol:    req.Symbol,
				Quantity:  req.Quantity,
				AveragePrice: req.Price,
			}
		} else {
			response.InternalError(w, "failed to get position")
			return
		}
	} else {
		if req.Side == "BUY" {
			totalCost := (pos.Quantity * pos.AveragePrice) + (req.Quantity * req.Price)
			pos.Quantity += req.Quantity
			pos.AveragePrice = totalCost / pos.Quantity
		} else {
			pos.Quantity -= req.Quantity
			if pos.Quantity <= 0 {
				pos.Quantity = 0
				pos.AveragePrice = 0
			}
		}
	}

	if err := h.posRepo.Upsert(pos); err != nil {
		response.InternalError(w, "failed to update position")
		return
	}

	response.Success(w, toPositionResponse(pos))
}

func toPositionResponse(pos *models.Position) PositionResponse {
	costBasis := pos.Quantity * pos.AveragePrice
	marketValue := costBasis
	unrealizedPnL := marketValue - costBasis

	return PositionResponse{
		ID:            pos.ID.String(),
		Symbol:        pos.Symbol,
		Quantity:      pos.Quantity,
		AveragePrice:  pos.AveragePrice,
		CurrentPrice:  pos.AveragePrice,
		UnrealizedPnL: unrealizedPnL,
		CostBasis:     costBasis,
		MarketValue:   marketValue,
		UpdatedAt:     pos.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
