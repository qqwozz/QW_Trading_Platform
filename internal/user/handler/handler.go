package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/qw-trading/platform/internal/models"
	"github.com/qw-trading/platform/internal/user/repository"
	"github.com/qw-trading/platform/pkg/errors"
	"github.com/qw-trading/platform/pkg/middleware"
	"github.com/qw-trading/platform/pkg/response"
)

type Handler struct {
	repo    *repository.UserRepository
	jwtSecret string
	jwtExpiry int
}

func New(repo *repository.UserRepository, jwtSecret string, jwtExpiry int) *Handler {
	return &Handler{repo: repo, jwtSecret: jwtSecret, jwtExpiry: jwtExpiry}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if req.Email == "" || req.Username == "" || req.Password == "" {
		response.BadRequest(w, "email, username, and password are required")
		return
	}

	exists, err := h.repo.EmailExists(req.Email)
	if err != nil {
		response.InternalError(w, "failed to check email")
		return
	}
	if exists {
		response.Conflict(w, "email already exists")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.InternalError(w, "failed to hash password")
		return
	}

	user := &models.User{
		ID:           uuid.New(),
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: string(hash),
		Status:       models.UserStatusActive,
	}

	if err := h.repo.Create(user); err != nil {
		response.InternalError(w, "failed to create user")
		return
	}

	response.Created(w, UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		response.BadRequest(w, "email and password are required")
		return
	}

	user, err := h.repo.GetByEmail(req.Email)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok && appErr.Code == http.StatusNotFound {
			response.Unauthorized(w, "invalid credentials")
			return
		}
		response.InternalError(w, "failed to get user")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		response.Unauthorized(w, "invalid credentials")
		return
	}

	accessToken, err := h.generateToken(user.ID, time.Duration(h.jwtExpiry)*time.Hour)
	if err != nil {
		response.InternalError(w, "failed to generate token")
		return
	}

	refreshToken, err := h.generateToken(user.ID, 7*24*time.Hour)
	if err != nil {
		response.InternalError(w, "failed to generate refresh token")
		return
	}

	response.Success(w, AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    h.jwtExpiry * 3600,
	})
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Unauthorized(w, "unauthorized")
		return
	}

	user, err := h.repo.GetByID(userID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok && appErr.Code == http.StatusNotFound {
			response.NotFound(w, "user not found")
			return
		}
		response.InternalError(w, "failed to get user")
		return
	}

	response.Success(w, UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	})
}

func (h *Handler) generateToken(userID uuid.UUID, expiry time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(expiry).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}
