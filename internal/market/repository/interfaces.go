package repository

import (
	"github.com/qw-trading/platform/internal/models"
)

type MarketRepositoryInterface interface {
	UpsertTicker(ticker *models.MarketTicker) error
	GetTicker(symbol string) (*models.MarketTicker, error)
	GetTickers() ([]models.MarketTicker, error)
	SaveOrderBookSnapshot(symbol string, bids, asks []models.OrderBookLevel) error
	GetRecentSnapshot(symbol string) (*models.OrderBookSnapshot, error)
}
