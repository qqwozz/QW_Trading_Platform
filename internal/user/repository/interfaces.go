package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/qw-trading/platform/internal/models"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	EmailExists(ctx context.Context, email string) (bool, error)
}
