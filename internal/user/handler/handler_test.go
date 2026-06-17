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
	apperr "github.com/qw-trading/platform/pkg/errors"
	"github.com/qw-trading/platform/pkg/middleware"
)

func TestRegister_Success(t *testing.T) {
	repo := &mockUserRepo{
		emailExistsFn: func(email string) (bool, error) { return false, nil },
		createFn:      func(user *models.User) error { return nil },
	}
	h := New(repo, "secret", 1)

	body, _ := json.Marshal(RegisterRequest{
		Email:    "test@example.com",
		Username: "tester",
		Password: "pass123",
	})
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Register(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", w.Code, http.StatusCreated)
	}
}

func TestRegister_EmptyFields(t *testing.T) {
	repo := &mockUserRepo{}
	h := New(repo, "secret", 1)

	body, _ := json.Marshal(RegisterRequest{Email: "test@example.com"})
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestRegister_EmailExists(t *testing.T) {
	repo := &mockUserRepo{
		emailExistsFn: func(email string) (bool, error) { return true, nil },
	}
	h := New(repo, "secret", 1)

	body, _ := json.Marshal(RegisterRequest{
		Email:    "exists@example.com",
		Username: "tester",
		Password: "pass123",
	})
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Register(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("status = %d, want %d", w.Code, http.StatusConflict)
	}
}

func TestRegister_InvalidJSON(t *testing.T) {
	repo := &mockUserRepo{}
	h := New(repo, "secret", 1)

	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader([]byte("not json")))
	w := httptest.NewRecorder()
	h.Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestLogin_Success(t *testing.T) {
	repo := &mockUserRepo{
		getByEmailFn: func(email string) (*models.User, error) {
			hash := "$2a$10$abcdefghijklmnopqrstuuABCDEFGHIJKLMNOPQRSTUVWXYZ01"
			return &models.User{
				ID:           uuid.New(),
				Email:        email,
				PasswordHash: hash,
				Status:       models.UserStatusActive,
			}, nil
		},
	}
	h := New(repo, "secret", 1)

	// Use bcrypt-hashed "password" for the test
	body, _ := json.Marshal(LoginRequest{Email: "test@example.com", Password: "wrong"})
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Login(w, req)

	// Should fail with wrong password
	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d (wrong password)", w.Code, http.StatusUnauthorized)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := &mockUserRepo{
		getByEmailFn: func(email string) (*models.User, error) {
			return nil, apperr.NotFound("user not found")
		},
	}
	h := New(repo, "secret", 1)

	body, _ := json.Marshal(LoginRequest{Email: "nobody@example.com", Password: "pass"})
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Login(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestLogin_EmptyFields(t *testing.T) {
	repo := &mockUserRepo{}
	h := New(repo, "secret", 1)

	body, _ := json.Marshal(LoginRequest{Email: ""})
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Login(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGetProfile_Success(t *testing.T) {
	userID := uuid.New()
	repo := &mockUserRepo{
		getByIDFn: func(id uuid.UUID) (*models.User, error) {
			return &models.User{
				ID:        id,
				Email:     "test@example.com",
				Username:  "tester",
				CreatedAt: time.Now(),
				Status:    models.UserStatusActive,
			}, nil
		},
	}
	h := New(repo, "secret", 1)

	req := httptest.NewRequest("GET", "/users/me", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID.String())
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	h.GetProfile(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetProfile_Unauthorized(t *testing.T) {
	repo := &mockUserRepo{}
	h := New(repo, "secret", 1)

	req := httptest.NewRequest("GET", "/users/me", nil)
	w := httptest.NewRecorder()
	h.GetProfile(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}
