package repository

import (
	"github.com/google/uuid"
	"github.com/qw-trading/platform/internal/models"
)

type UserRepositoryInterface interface {
	Create(user *models.User) error
	GetByID(id uuid.UUID) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	EmailExists(email string) (bool, error)
}
