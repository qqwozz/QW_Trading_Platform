package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/qw-trading/platform/internal/models"
	apperr "github.com/qw-trading/platform/pkg/errors"
	"github.com/qw-trading/platform/pkg/middleware"
)

func TestListAccounts_Success(t *testing.T) {
	userID := uuid.New()
	repo := &mockAccountRepo{
		getByUserIDFn: func(uid uuid.UUID) ([]models.Account, error) {
			return []models.Account{
				{ID: uuid.New(), UserID: uid, Type: models.AccountTypeCash, Balance: 1000, Currency: "USDT", Status: models.AccountStatusActive},
			}, nil
		},
	}
	h := New(repo)

	req := httptest.NewRequest("GET", "/accounts", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID.String())
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	h.ListAccounts(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestListAccounts_Unauthorized(t *testing.T) {
	repo := &mockAccountRepo{}
	h := New(repo)

	req := httptest.NewRequest("GET", "/accounts", nil)
	w := httptest.NewRecorder()
	h.ListAccounts(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestDeposit_Success(t *testing.T) {
	userID := uuid.New()
	accountID := uuid.New()
	repo := &mockAccountRepo{
		getByUserAndCurrencyFn: func(uid uuid.UUID, currency string) (*models.Account, error) {
			return &models.Account{ID: accountID, UserID: uid, Balance: 500, Currency: "USDT"}, nil
		},
		creditFn: func(id uuid.UUID, amount float64) error { return nil },
	}
	h := New(repo)

	body, _ := json.Marshal(DepositRequest{Currency: "USDT", Amount: 100})
	req := httptest.NewRequest("POST", "/accounts/deposit", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID.String())
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	h.Deposit(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestDeposit_InvalidAmount(t *testing.T) {
	repo := &mockAccountRepo{}
	h := New(repo)

	body, _ := json.Marshal(DepositRequest{Currency: "USDT", Amount: -10})
	req := httptest.NewRequest("POST", "/accounts/deposit", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, uuid.New().String())
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	h.Deposit(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestDeposit_AccountNotFound(t *testing.T) {
	repo := &mockAccountRepo{
		getByUserAndCurrencyFn: func(uid uuid.UUID, currency string) (*models.Account, error) {
			return nil, apperr.NotFound("account not found")
		},
	}
	h := New(repo)

	body, _ := json.Marshal(DepositRequest{Currency: "USDT", Amount: 100})
	req := httptest.NewRequest("POST", "/accounts/deposit", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, uuid.New().String())
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	h.Deposit(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestWithdraw_Success(t *testing.T) {
	userID := uuid.New()
	accountID := uuid.New()
	repo := &mockAccountRepo{
		getByUserAndCurrencyFn: func(uid uuid.UUID, currency string) (*models.Account, error) {
			return &models.Account{ID: accountID, UserID: uid, Balance: 1000, FrozenBalance: 0, Currency: "USDT"}, nil
		},
		debitFn: func(id uuid.UUID, amount float64) error { return nil },
	}
	h := New(repo)

	body, _ := json.Marshal(WithdrawRequest{Currency: "USDT", Amount: 200})
	req := httptest.NewRequest("POST", "/accounts/withdraw", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID.String())
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	h.Withdraw(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestWithdraw_InsufficientBalance(t *testing.T) {
	repo := &mockAccountRepo{
		getByUserAndCurrencyFn: func(uid uuid.UUID, currency string) (*models.Account, error) {
			return &models.Account{Balance: 100, FrozenBalance: 0, Currency: "USDT"}, nil
		},
	}
	h := New(repo)

	body, _ := json.Marshal(WithdrawRequest{Currency: "USDT", Amount: 500})
	req := httptest.NewRequest("POST", "/accounts/withdraw", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, uuid.New().String())
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	h.Withdraw(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d (insufficient balance)", w.Code, http.StatusBadRequest)
	}
}

func TestWithdraw_InvalidAmount(t *testing.T) {
	repo := &mockAccountRepo{}
	h := New(repo)

	body, _ := json.Marshal(WithdrawRequest{Currency: "USDT", Amount: -10})
	req := httptest.NewRequest("POST", "/accounts/withdraw", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, uuid.New().String())
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	h.Withdraw(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}
