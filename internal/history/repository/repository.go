package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/qw-trading/platform/internal/db"
	"github.com/qw-trading/platform/internal/models"
	apperr "github.com/qw-trading/platform/pkg/errors"
)

type HistoryRepository struct {
	db *db.Database
}

func New(database *db.Database) *HistoryRepository {
	return &HistoryRepository{db: database}
}

func (r *HistoryRepository) RecordBalanceChange(ctx context.Context, entry *models.BalanceHistory) error {
	query := `
		INSERT INTO balance_history (id, user_id, account_id, currency, amount, balance_before, balance_after, type, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at`

	return r.db.QueryRowContext(ctx, query,
		entry.ID, entry.UserID, entry.AccountID, entry.Currency,
		entry.Amount, entry.BalanceBefore, entry.BalanceAfter,
		entry.Type, entry.Description,
	).Scan(&entry.CreatedAt)
}

func (r *HistoryRepository) RecordPositionChange(ctx context.Context, entry *models.PositionHistory) error {
	query := `
		INSERT INTO position_history (id, user_id, account_id, symbol, quantity_change, quantity_before, quantity_after, avg_price_before, avg_price_after, type, trade_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING created_at`

	return r.db.QueryRowContext(ctx, query,
		entry.ID, entry.UserID, entry.AccountID, entry.Symbol,
		entry.QuantityChange, entry.QuantityBefore, entry.QuantityAfter,
		entry.AvgPriceBefore, entry.AvgPriceAfter, entry.Type, entry.TradeID,
	).Scan(&entry.CreatedAt)
}

func (r *HistoryRepository) GetOrderHistory(ctx context.Context, userID uuid.UUID, symbol string, limit, offset int) ([]models.Order, int, error) {
	countQuery := `SELECT COUNT(*) FROM orders WHERE user_id = $1`
	listQuery := `
		SELECT id, user_id, account_id, symbol, side, type, price, quantity,
		       filled_quantity, status, time_in_force, created_at, updated_at, expired_at
		FROM orders WHERE user_id = $1`

	args := []interface{}{userID}
	argIdx := 1

	if symbol != "" {
		argIdx++
		countQuery += fmt.Sprintf(` AND symbol = $%d`, argIdx)
		listQuery += fmt.Sprintf(` AND symbol = $%d`, argIdx)
		args = append(args, symbol)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args[:argIdx]...).Scan(&total); err != nil {
		return nil, 0, apperr.InternalErr("failed to count orders", err)
	}

	listQuery += fmt.Sprintf(` ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, argIdx+1, argIdx+2)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, apperr.InternalErr("failed to list orders", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(
			&order.ID, &order.UserID, &order.AccountID, &order.Symbol,
			&order.Side, &order.Type, &order.Price, &order.Quantity,
			&order.FilledQuantity, &order.Status, &order.TimeInForce,
			&order.CreatedAt, &order.UpdatedAt, &order.ExpiredAt,
		); err != nil {
			return nil, 0, apperr.InternalErr("failed to scan order", err)
		}
		orders = append(orders, order)
	}
	return orders, total, nil
}

func (r *HistoryRepository) GetTradeHistory(ctx context.Context, userID uuid.UUID, symbol string, limit, offset int) ([]models.Trade, int, error) {
	countQuery := `SELECT COUNT(*) FROM trades WHERE buyer_id = $1 OR seller_id = $1`
	listQuery := `
		SELECT id, symbol, buyer_order_id, seller_order_id, buyer_id, seller_id,
		       price, quantity, buyer_fee, seller_fee, executed_at
		FROM trades WHERE buyer_id = $1 OR seller_id = $1`

	args := []interface{}{userID}
	argIdx := 1

	if symbol != "" {
		argIdx++
		countQuery += fmt.Sprintf(` AND symbol = $%d`, argIdx)
		listQuery += fmt.Sprintf(` AND symbol = $%d`, argIdx)
		args = append(args, symbol)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args[:argIdx]...).Scan(&total); err != nil {
		return nil, 0, apperr.InternalErr("failed to count trades", err)
	}

	listQuery += fmt.Sprintf(` ORDER BY executed_at DESC LIMIT $%d OFFSET $%d`, argIdx+1, argIdx+2)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, apperr.InternalErr("failed to list trades", err)
	}
	defer rows.Close()

	var trades []models.Trade
	for rows.Next() {
		var trade models.Trade
		if err := rows.Scan(
			&trade.ID, &trade.Symbol, &trade.BuyerOrderID, &trade.SellerOrderID,
			&trade.BuyerID, &trade.SellerID, &trade.Price, &trade.Quantity,
			&trade.BuyerFee, &trade.SellerFee, &trade.ExecutedAt,
		); err != nil {
			return nil, 0, apperr.InternalErr("failed to scan trade", err)
		}
		trades = append(trades, trade)
	}
	return trades, total, nil
}

func (r *HistoryRepository) GetBalanceHistory(ctx context.Context, userID uuid.UUID, currency string, limit, offset int) ([]models.BalanceHistory, int, error) {
	countQuery := `SELECT COUNT(*) FROM balance_history WHERE user_id = $1`
	listQuery := `
		SELECT id, user_id, account_id, currency, amount, balance_before, balance_after, type, description, created_at
		FROM balance_history WHERE user_id = $1`

	args := []interface{}{userID}
	argIdx := 1

	if currency != "" {
		argIdx++
		countQuery += fmt.Sprintf(` AND currency = $%d`, argIdx)
		listQuery += fmt.Sprintf(` AND currency = $%d`, argIdx)
		args = append(args, currency)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args[:argIdx]...).Scan(&total); err != nil {
		return nil, 0, apperr.InternalErr("failed to count balance history", err)
	}

	listQuery += fmt.Sprintf(` ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, argIdx+1, argIdx+2)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, apperr.InternalErr("failed to list balance history", err)
	}
	defer rows.Close()

	var history []models.BalanceHistory
	for rows.Next() {
		var entry models.BalanceHistory
		if err := rows.Scan(
			&entry.ID, &entry.UserID, &entry.AccountID, &entry.Currency,
			&entry.Amount, &entry.BalanceBefore, &entry.BalanceAfter,
			&entry.Type, &entry.Description, &entry.CreatedAt,
		); err != nil {
			return nil, 0, apperr.InternalErr("failed to scan balance history", err)
		}
		history = append(history, entry)
	}
	return history, total, nil
}

func (r *HistoryRepository) GetPositionHistory(ctx context.Context, userID uuid.UUID, symbol string, limit, offset int) ([]models.PositionHistory, int, error) {
	countQuery := `SELECT COUNT(*) FROM position_history WHERE user_id = $1`
	listQuery := `
		SELECT id, user_id, account_id, symbol, quantity_change, quantity_before, quantity_after,
		       avg_price_before, avg_price_after, type, trade_id, created_at
		FROM position_history WHERE user_id = $1`

	args := []interface{}{userID}
	argIdx := 1

	if symbol != "" {
		argIdx++
		countQuery += fmt.Sprintf(` AND symbol = $%d`, argIdx)
		listQuery += fmt.Sprintf(` AND symbol = $%d`, argIdx)
		args = append(args, symbol)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args[:argIdx]...).Scan(&total); err != nil {
		return nil, 0, apperr.InternalErr("failed to count position history", err)
	}

	listQuery += fmt.Sprintf(` ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, argIdx+1, argIdx+2)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, apperr.InternalErr("failed to list position history", err)
	}
	defer rows.Close()

	var history []models.PositionHistory
	for rows.Next() {
		var entry models.PositionHistory
		if err := rows.Scan(
			&entry.ID, &entry.UserID, &entry.AccountID, &entry.Symbol,
			&entry.QuantityChange, &entry.QuantityBefore, &entry.QuantityAfter,
			&entry.AvgPriceBefore, &entry.AvgPriceAfter, &entry.Type,
			&entry.TradeID, &entry.CreatedAt,
		); err != nil {
			return nil, 0, apperr.InternalErr("failed to scan position history", err)
		}
		history = append(history, entry)
	}
	return history, total, nil
}
