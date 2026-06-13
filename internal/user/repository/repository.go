// Package repository provides data access for the user domain, backed by PostgreSQL.
package repository

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/qw-trading/platform/internal/db"
	"github.com/qw-trading/platform/internal/models"
	apperr "github.com/qw-trading/platform/pkg/errors"
)

// UserRepository handles database operations for user entities.
type UserRepository struct {
	db *db.Database
}

// New creates a new UserRepository.
func New(db *db.Database) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user and populates the CreatedAt and UpdatedAt fields.
func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (id, email, username, password_hash, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, updated_at`

	return r.db.QueryRow(query,
		user.ID, user.Email, user.Username, user.PasswordHash, user.Status,
	).Scan(&user.CreatedAt, &user.UpdatedAt)
}

// GetByID retrieves a user by their UUID. Returns a NotFound error if no
// matching user exists.
func (r *UserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, username, password_hash, created_at, updated_at, status
		FROM users WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.CreatedAt, &user.UpdatedAt, &user.Status,
	)
	if err == sql.ErrNoRows {
		return nil, apperr.NotFound("user not found")
	}
	return user, apperr.InternalErr("failed to get user", err)
}

// GetByEmail retrieves a user by their email address. Returns a NotFound
// error if no matching user exists.
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, username, password_hash, created_at, updated_at, status
		FROM users WHERE email = $1`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.CreatedAt, &user.UpdatedAt, &user.Status,
	)
	if err == sql.ErrNoRows {
		return nil, apperr.NotFound("user not found")
	}
	return user, apperr.InternalErr("failed to get user", err)
}

// EmailExists checks whether a user with the given email address exists.
func (r *UserRepository) EmailExists(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := r.db.QueryRow(query, email).Scan(&exists)
	return exists, err
}
