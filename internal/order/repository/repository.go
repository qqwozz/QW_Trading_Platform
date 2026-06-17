package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/qw-trading/platform/internal/db"
	"github.com/qw-trading/platform/internal/models"
	apperr "github.com/qw-trading/platform/pkg/errors"
)

type OrderRepository struct {
	db *db.Database
}

func New(db *db.Database) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, order *models.Order) error {
	query := `
		INSERT INTO orders (id, user_id, account_id, symbol, side, type, price, quantity, filled_quantity, status, time_in_force)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING created_at, updated_at`

	return r.db.QueryRowContext(ctx, query,
		order.ID, order.UserID, order.AccountID, order.Symbol, order.Side,
		order.Type, order.Price, order.Quantity, order.FilledQuantity,
		order.Status, order.TimeInForce,
	).Scan(&order.CreatedAt, &order.UpdatedAt)
}

func (r *OrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	order := &models.Order{}
	query := `
		SELECT id, user_id, account_id, symbol, side, type, price, quantity, 
		       filled_quantity, status, time_in_force, created_at, updated_at
		FROM orders WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&order.ID, &order.UserID, &order.AccountID, &order.Symbol,
		&order.Side, &order.Type, &order.Price, &order.Quantity,
		&order.FilledQuantity, &order.Status, &order.TimeInForce,
		&order.CreatedAt, &order.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, apperr.NotFound("order not found")
	}
	return order, apperr.InternalErr("failed to get order", err)
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus, filledQuantity float64) error {
	query := `
		UPDATE orders SET status = $2, filled_quantity = $3, updated_at = NOW()
		WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, orderID, status, filledQuantity)
	if err != nil {
		return apperr.InternalErr("failed to update order", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apperr.NotFound("order not found")
	}
	return nil
}

type ListFilter struct {
	UserID uuid.UUID
	Symbol string
	Status string
	Limit  int
	Offset int
}

func (r *OrderRepository) List(ctx context.Context, filter ListFilter) ([]models.Order, int, error) {
	countQuery := `SELECT COUNT(*) FROM orders WHERE user_id = $1`
	listQuery := `
		SELECT id, user_id, account_id, symbol, side, type, price, quantity, 
		       filled_quantity, status, time_in_force, created_at, updated_at
		FROM orders WHERE user_id = $1`

	args := []interface{}{filter.UserID}
	argIdx := 1

	if filter.Symbol != "" {
		argIdx++
		countQuery += fmt.Sprintf(` AND symbol = $%d`, argIdx)
		listQuery += fmt.Sprintf(` AND symbol = $%d`, argIdx)
		args = append(args, filter.Symbol)
	}

	if filter.Status != "" {
		argIdx++
		countQuery += fmt.Sprintf(` AND status = $%d`, argIdx)
		listQuery += fmt.Sprintf(` AND status = $%d`, argIdx)
		args = append(args, filter.Status)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args[:argIdx]...).Scan(&total); err != nil {
		return nil, 0, apperr.InternalErr("failed to count orders", err)
	}

	listQuery += fmt.Sprintf(` ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, argIdx+1, argIdx+2)
	args = append(args, filter.Limit, filter.Offset)

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
			&order.CreatedAt, &order.UpdatedAt,
		); err != nil {
			return nil, 0, apperr.InternalErr("failed to scan order", err)
		}
		orders = append(orders, order)
	}
	return orders, total, nil
}

type TradeRepository struct {
	db *db.Database
}

func NewTradeRepository(db *db.Database) *TradeRepository {
	return &TradeRepository{db: db}
}

func (r *TradeRepository) Create(ctx context.Context, trade *models.Trade) error {
	query := `
		INSERT INTO trades (id, symbol, buyer_order_id, seller_order_id, buyer_id, seller_id, price, quantity, buyer_fee, seller_fee)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING executed_at`

	return r.db.QueryRowContext(ctx, query,
		trade.ID, trade.Symbol, trade.BuyerOrderID, trade.SellerOrderID,
		trade.BuyerID, trade.SellerID, trade.Price, trade.Quantity,
		trade.BuyerFee, trade.SellerFee,
	).Scan(&trade.ExecutedAt)
}
