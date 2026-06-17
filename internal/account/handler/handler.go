package handler

import (
	"encoding/json"
	"net/http"

	"github.com/qw-trading/platform/internal/account/repository"
	"github.com/qw-trading/platform/pkg/middleware"
	"github.com/qw-trading/platform/pkg/response"
)

type Handler struct {
	repo repository.AccountRepositoryInterface
}

func New(repo repository.AccountRepositoryInterface) *Handler {
	return &Handler{repo: repo}
}

type AccountResponse struct {
	ID            string  `json:"id"`
	UserID        string  `json:"user_id"`
	Type          string  `json:"type"`
	Balance       float64 `json:"balance"`
	FrozenBalance float64 `json:"frozen_balance"`
	Currency      string  `json:"currency"`
	Status        string  `json:"status"`
}

type DepositRequest struct {
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}

type DepositResponse struct {
	AccountID  string  `json:"account_id"`
	NewBalance float64 `json:"new_balance"`
}

func (h *Handler) ListAccounts(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Unauthorized(w, "unauthorized")
		return
	}

	accounts, err := h.repo.GetByUserID(r.Context(), userID)
	if err != nil {
		response.InternalError(w, "failed to get accounts")
		return
	}

	result := make([]AccountResponse, len(accounts))
	for i, acc := range accounts {
		result[i] = AccountResponse{
			ID:            acc.ID.String(),
			UserID:        acc.UserID.String(),
			Type:          string(acc.Type),
			Balance:       acc.Balance,
			FrozenBalance: acc.FrozenBalance,
			Currency:      acc.Currency,
			Status:        string(acc.Status),
		}
	}

	response.Success(w, map[string]interface{}{"accounts": result})
}

func (h *Handler) Deposit(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Unauthorized(w, "unauthorized")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req DepositRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if req.Currency == "" || req.Amount <= 0 || req.Amount > 10_000_000 {
		response.BadRequest(w, "currency and amount (0, 10000000] are required")
		return
	}

	account, err := h.repo.GetByUserAndCurrency(r.Context(), userID, req.Currency)
	if err != nil {
		response.NotFound(w, "account not found")
		return
	}

	if err := h.repo.Credit(r.Context(), account.ID, req.Amount); err != nil {
		response.InternalError(w, "failed to deposit")
		return
	}

	response.Success(w, DepositResponse{
		AccountID:  account.ID.String(),
		NewBalance: account.Balance + req.Amount,
	})
}
