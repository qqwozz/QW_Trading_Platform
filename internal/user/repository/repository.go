package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/qw-trading/platform/internal/db"
	"github.com/qw-trading/platform/internal/models"
	apperr "github.com/qw-trading/platform/pkg/errors"
)

type UserRepository struct {
	db *db.Database
}

func New(db *db.Database) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, email, username, password_hash, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, updated_at`

	return r.db.QueryRowContext(ctx, query,
		user.ID, user.Email, user.Username, user.PasswordHash, user.Status,
	).Scan(&user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, username, password_hash, created_at, updated_at, status
		FROM users WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.CreatedAt, &user.UpdatedAt, &user.Status,
	)
	if err == sql.ErrNoRows {
		return nil, apperr.NotFound("user not found")
	}
	return user, apperr.InternalErr("failed to get user", err)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, username, password_hash, created_at, updated_at, status
		FROM users WHERE email = $1`

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.CreatedAt, &user.UpdatedAt, &user.Status,
	)
	if err == sql.ErrNoRows {
		return nil, apperr.NotFound("user not found")
	}
	return user, apperr.InternalErr("failed to get user", err)
}

func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	return exists, err
}
