package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/qw-trading/platform/internal/models"
)

type PositionRepositoryInterface interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.Position, error)
	GetByUserAndSymbol(ctx context.Context, userID uuid.UUID, symbol string) (*models.Position, error)
	Upsert(ctx context.Context, pos *models.Position) error
}

type AccountRepositoryInterface interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.Account, error)
}
