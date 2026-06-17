// Package handler implements HTTP handlers for account management endpoints.
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/qw-trading/platform/internal/account/repository"
	"github.com/qw-trading/platform/pkg/middleware"
	"github.com/qw-trading/platform/pkg/response"
)

// Handler holds dependencies for account-related HTTP handlers.
type Handler struct {
	repo repository.AccountRepositoryInterface
}

// New creates a new Handler with the given repository.
func New(repo repository.AccountRepositoryInterface) *Handler {
	return &Handler{repo: repo}
}

// AccountResponse is the JSON response containing account information.
type AccountResponse struct {
	ID            string  `json:"id"`
	UserID        string  `json:"user_id"`
	Type          string  `json:"type"`
	Balance       float64 `json:"balance"`
	FrozenBalance float64 `json:"frozen_balance"`
	Currency      string  `json:"currency"`
	Status        string  `json:"status"`
}

// DepositRequest is the JSON request body for depositing funds.
type DepositRequest struct {
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}

// DepositResponse is the JSON response after a successful deposit.
type DepositResponse struct {
	AccountID  string  `json:"account_id"`
	NewBalance float64 `json:"new_balance"`
}

// ListAccounts handles GET /accounts. Returns all accounts for the authenticated user.
func (h *Handler) ListAccounts(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Unauthorized(w, "unauthorized")
		return
	}

	accounts, err := h.repo.GetByUserID(userID)
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

// Deposit handles POST /accounts/deposit. It credits the specified amount to
// the user's account for the given currency.
func (h *Handler) Deposit(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Unauthorized(w, "unauthorized")
		return
	}

	var req DepositRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if req.Currency == "" || req.Amount <= 0 {
		response.BadRequest(w, "currency and positive amount are required")
		return
	}

	account, err := h.repo.GetByUserAndCurrency(userID, req.Currency)
	if err != nil {
		response.NotFound(w, "account not found")
		return
	}

	if err := h.repo.Credit(account.ID, req.Amount); err != nil {
		response.InternalError(w, "failed to deposit")
		return
	}

	response.Success(w, DepositResponse{
		AccountID:  account.ID.String(),
		NewBalance: account.Balance + req.Amount,
	})
}
