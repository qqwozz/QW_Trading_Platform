package repository

import (
	"context"

	"github.com/qw-trading/platform/internal/models"
)

type MarketRepositoryInterface interface {
	UpsertTicker(ctx context.Context, ticker *models.MarketTicker) error
	GetTicker(ctx context.Context, symbol string) (*models.MarketTicker, error)
	GetTickers(ctx context.Context) ([]models.MarketTicker, error)
	SaveOrderBookSnapshot(ctx context.Context, symbol string, bids, asks []models.OrderBookLevel) error
	GetRecentSnapshot(ctx context.Context, symbol string) (*models.OrderBookSnapshot, error)
}
