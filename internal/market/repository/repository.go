package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/qw-trading/platform/internal/db"
	"github.com/qw-trading/platform/internal/models"
	apperr "github.com/qw-trading/platform/pkg/errors"
)

type MarketRepository struct {
	db *db.Database
}

func New(database *db.Database) *MarketRepository {
	return &MarketRepository{db: database}
}

func (r *MarketRepository) UpsertTicker(ctx context.Context, ticker *models.MarketTicker) error {
	query := `
		INSERT INTO market_tickers (id, symbol, last_price, best_bid, best_ask, bid_volume, ask_volume, volume_24h, high_24h, low_24h, change_24h, change_pct_24h, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (symbol) DO UPDATE SET
			last_price = EXCLUDED.last_price,
			best_bid = EXCLUDED.best_bid,
			best_ask = EXCLUDED.best_ask,
			bid_volume = EXCLUDED.bid_volume,
			ask_volume = EXCLUDED.ask_volume,
			volume_24h = EXCLUDED.volume_24h,
			high_24h = EXCLUDED.high_24h,
			low_24h = EXCLUDED.low_24h,
			change_24h = EXCLUDED.change_24h,
			change_pct_24h = EXCLUDED.change_pct_24h,
			updated_at = EXCLUDED.updated_at`

	_, err := r.db.ExecContext(ctx, query,
		ticker.ID, ticker.Symbol, ticker.LastPrice, ticker.BestBid, ticker.BestAsk,
		ticker.BidVolume, ticker.AskVolume, ticker.Volume24h, ticker.High24h,
		ticker.Low24h, ticker.Change24h, ticker.ChangePct24h, ticker.UpdatedAt,
	)
	return err
}

func (r *MarketRepository) GetTicker(ctx context.Context, symbol string) (*models.MarketTicker, error) {
	ticker := &models.MarketTicker{}
	query := `
		SELECT id, symbol, last_price, best_bid, best_ask, bid_volume, ask_volume,
			volume_24h, high_24h, low_24h, change_24h, change_pct_24h, updated_at
		FROM market_tickers WHERE symbol = $1`

	err := r.db.QueryRowContext(ctx, query, symbol).Scan(
		&ticker.ID, &ticker.Symbol, &ticker.LastPrice, &ticker.BestBid, &ticker.BestAsk,
		&ticker.BidVolume, &ticker.AskVolume, &ticker.Volume24h, &ticker.High24h,
		&ticker.Low24h, &ticker.Change24h, &ticker.ChangePct24h, &ticker.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, apperr.NotFound("ticker not found")
	}
	if err != nil {
		return nil, apperr.InternalErr("failed to get ticker", err)
	}
	return ticker, nil
}

func (r *MarketRepository) GetTickers(ctx context.Context) ([]models.MarketTicker, error) {
	query := `
		SELECT id, symbol, last_price, best_bid, best_ask, bid_volume, ask_volume,
			volume_24h, high_24h, low_24h, change_24h, change_pct_24h, updated_at
		FROM market_tickers ORDER BY symbol`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, apperr.InternalErr("failed to get tickers", err)
	}
	defer rows.Close()

	tickers := make([]models.MarketTicker, 0)
	for rows.Next() {
		var t models.MarketTicker
		if err := rows.Scan(
			&t.ID, &t.Symbol, &t.LastPrice, &t.BestBid, &t.BestAsk,
			&t.BidVolume, &t.AskVolume, &t.Volume24h, &t.High24h,
			&t.Low24h, &t.Change24h, &t.ChangePct24h, &t.UpdatedAt,
		); err != nil {
			return nil, apperr.InternalErr("failed to scan ticker", err)
		}
		tickers = append(tickers, t)
	}
	if err := rows.Err(); err != nil {
		return nil, apperr.InternalErr("failed to iterate tickers", err)
	}
	return tickers, nil
}

func (r *MarketRepository) SaveOrderBookSnapshot(ctx context.Context, symbol string, bids, asks []models.OrderBookLevel) error {
	bidsJSON, err := json.Marshal(bids)
	if err != nil {
		return apperr.InternalErr("failed to marshal bids", err)
	}
	asksJSON, err := json.Marshal(asks)
	if err != nil {
		return apperr.InternalErr("failed to marshal asks", err)
	}

	query := `
		INSERT INTO order_book_snapshots (id, symbol, bids_json, asks_json, created_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4)`

	_, err = r.db.ExecContext(ctx, query, symbol, bidsJSON, asksJSON, time.Now().UTC())
	return err
}

func (r *MarketRepository) GetRecentSnapshot(ctx context.Context, symbol string) (*models.OrderBookSnapshot, error) {
	snapshot := &models.OrderBookSnapshot{}
	query := `
		SELECT id, symbol, bids_json, asks_json, created_at
		FROM order_book_snapshots
		WHERE symbol = $1
		ORDER BY created_at DESC
		LIMIT 1`

	var bidsJSON, asksJSON []byte
	err := r.db.QueryRowContext(ctx, query, symbol).Scan(
		&snapshot.ID, &snapshot.Symbol, &bidsJSON, &asksJSON, &snapshot.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, apperr.NotFound("no order book snapshot found")
	}
	if err != nil {
		return nil, apperr.InternalErr("failed to get order book snapshot", err)
	}

	if err := json.Unmarshal(bidsJSON, &snapshot.Bids); err != nil {
		return nil, apperr.InternalErr("failed to unmarshal bids", err)
	}
	if err := json.Unmarshal(asksJSON, &snapshot.Asks); err != nil {
		return nil, apperr.InternalErr("failed to unmarshal asks", err)
	}

	return snapshot, nil
}
