package repository

import (
	"github.com/google/uuid"
	"github.com/qw-trading/platform/internal/models"
)

type PositionRepositoryInterface interface {
	GetByUserID(userID uuid.UUID) ([]models.Position, error)
	GetByUserAndSymbol(userID uuid.UUID, symbol string) (*models.Position, error)
	Upsert(pos *models.Position) error
}

type AccountRepositoryInterface interface {
	GetByUserID(userID uuid.UUID) ([]models.Account, error)
}
