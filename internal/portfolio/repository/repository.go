package repository

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/qw-trading/platform/internal/db"
	"github.com/qw-trading/platform/internal/models"
	apperr "github.com/qw-trading/platform/pkg/errors"
)

type PositionRepository struct {
	db *db.Database
}

func New(db *db.Database) *PositionRepository {
	return &PositionRepository{db: db}
}

func (r *PositionRepository) GetByUserID(userID uuid.UUID) ([]models.Position, error) {
	query := `
		SELECT id, user_id, account_id, symbol, quantity, average_price, unrealized_pnl, updated_at
		FROM positions WHERE user_id = $1 AND quantity > 0
		ORDER BY updated_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, apperr.InternalErr("failed to query positions", err)
	}
	defer rows.Close()

	var positions []models.Position
	for rows.Next() {
		var pos models.Position
		if err := rows.Scan(
			&pos.ID, &pos.UserID, &pos.AccountID, &pos.Symbol,
			&pos.Quantity, &pos.AveragePrice, &pos.UnrealizedPnL,
			&pos.UpdatedAt,
		); err != nil {
			return nil, apperr.InternalErr("failed to scan position", err)
		}
		positions = append(positions, pos)
	}
	return positions, nil
}

func (r *PositionRepository) GetByUserAndSymbol(userID uuid.UUID, symbol string) (*models.Position, error) {
	pos := &models.Position{}
	query := `
		SELECT id, user_id, account_id, symbol, quantity, average_price, unrealized_pnl, updated_at
		FROM positions WHERE user_id = $1 AND symbol = $2`

	err := r.db.QueryRow(query, userID, symbol).Scan(
		&pos.ID, &pos.UserID, &pos.AccountID, &pos.Symbol,
		&pos.Quantity, &pos.AveragePrice, &pos.UnrealizedPnL,
		&pos.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, apperr.NotFound("position not found")
	}
	return pos, apperr.InternalErr("failed to get position", err)
}

func (r *PositionRepository) Upsert(pos *models.Position) error {
	query := `
		INSERT INTO positions (id, user_id, account_id, symbol, quantity, average_price, unrealized_pnl)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id, symbol) 
		DO UPDATE SET 
			quantity = EXCLUDED.quantity,
			average_price = EXCLUDED.average_price,
			unrealized_pnl = EXCLUDED.unrealized_pnl,
			updated_at = NOW()
		RETURNING updated_at`

	return r.db.QueryRow(query,
		pos.ID, pos.UserID, pos.AccountID, pos.Symbol,
		pos.Quantity, pos.AveragePrice, pos.UnrealizedPnL,
	).Scan(&pos.UpdatedAt)
}

type AccountRepository struct {
	db *db.Database
}

func NewAccountRepository(db *db.Database) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) GetByUserID(userID uuid.UUID) ([]models.Account, error) {
	query := `
		SELECT id, user_id, type, balance, frozen_balance, currency, created_at, updated_at, status
		FROM accounts WHERE user_id = $1 ORDER BY created_at`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, apperr.InternalErr("failed to query accounts", err)
	}
	defer rows.Close()

	var accounts []models.Account
	for rows.Next() {
		var acc models.Account
		if err := rows.Scan(
			&acc.ID, &acc.UserID, &acc.Type, &acc.Balance,
			&acc.FrozenBalance, &acc.Currency, &acc.CreatedAt,
			&acc.UpdatedAt, &acc.Status,
		); err != nil {
			return nil, apperr.InternalErr("failed to scan account", err)
		}
		accounts = append(accounts, acc)
	}
	return accounts, nil
}
