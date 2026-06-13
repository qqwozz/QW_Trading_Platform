// Package repository provides data access for the account domain, backed by PostgreSQL.
package repository

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/qw-trading/platform/internal/db"
	"github.com/qw-trading/platform/internal/models"
	apperr "github.com/qw-trading/platform/pkg/errors"
)

// AccountRepository handles database operations for trading account entities.
type AccountRepository struct {
	db *db.Database
}

// New creates a new AccountRepository.
func New(db *db.Database) *AccountRepository {
	return &AccountRepository{db: db}
}

// Create inserts a new account and populates the CreatedAt and UpdatedAt fields.
func (r *AccountRepository) Create(account *models.Account) error {
	query := `
		INSERT INTO accounts (id, user_id, type, balance, frozen_balance, currency, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, updated_at`

	return r.db.QueryRow(query,
		account.ID, account.UserID, account.Type, account.Balance,
		account.FrozenBalance, account.Currency, account.Status,
	).Scan(&account.CreatedAt, &account.UpdatedAt)
}

// GetByID retrieves an account by its UUID. Returns a NotFound error if
// no matching account exists.
func (r *AccountRepository) GetByID(id uuid.UUID) (*models.Account, error) {
	account := &models.Account{}
	query := `
		SELECT id, user_id, type, balance, frozen_balance, currency, created_at, updated_at, status
		FROM accounts WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&account.ID, &account.UserID, &account.Type, &account.Balance,
		&account.FrozenBalance, &account.Currency, &account.CreatedAt,
		&account.UpdatedAt, &account.Status,
	)
	if err == sql.ErrNoRows {
		return nil, apperr.NotFound("account not found")
	}
	return account, apperr.InternalErr("failed to get account", err)
}

// GetByUserAndCurrency retrieves the account for a specific user and currency.
// Returns a NotFound error if no matching account exists.
func (r *AccountRepository) GetByUserAndCurrency(userID uuid.UUID, currency string) (*models.Account, error) {
	account := &models.Account{}
	query := `
		SELECT id, user_id, type, balance, frozen_balance, currency, created_at, updated_at, status
		FROM accounts WHERE user_id = $1 AND currency = $2`

	err := r.db.QueryRow(query, userID, currency).Scan(
		&account.ID, &account.UserID, &account.Type, &account.Balance,
		&account.FrozenBalance, &account.Currency, &account.CreatedAt,
		&account.UpdatedAt, &account.Status,
	)
	if err == sql.ErrNoRows {
		return nil, apperr.NotFound("account not found")
	}
	return account, apperr.InternalErr("failed to get account", err)
}

// GetByUserID retrieves all accounts for a user, ordered by creation time.
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
		var account models.Account
		if err := rows.Scan(
			&account.ID, &account.UserID, &account.Type, &account.Balance,
			&account.FrozenBalance, &account.Currency, &account.CreatedAt,
			&account.UpdatedAt, &account.Status,
		); err != nil {
			return nil, apperr.InternalErr("failed to scan account", err)
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

// FreezeBalance atomically transfers funds from available balance to frozen balance
// within a database transaction. Returns an error if the account is not found or
// has insufficient funds.
func (r *AccountRepository) FreezeBalance(accountID uuid.UUID, amount float64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return apperr.InternalErr("failed to begin transaction", err)
	}
	defer tx.Rollback()

	// Lock the row for update to prevent concurrent balance modifications.
	var balance float64
	err = tx.QueryRow(
		`SELECT balance FROM accounts WHERE id = $1 FOR UPDATE`,
		accountID,
	).Scan(&balance)
	if err != nil {
		return apperr.NotFound("account not found")
	}

	if balance < amount {
		return apperr.BadRequest("insufficient balance")
	}

	_, err = tx.Exec(
		`UPDATE accounts SET balance = balance - $2, frozen_balance = frozen_balance + $2, updated_at = NOW() WHERE id = $1`,
		accountID, amount,
	)
	if err != nil {
		return apperr.InternalErr("failed to freeze balance", err)
	}

	return tx.Commit()
}

// UnfreezeBalance atomically transfers funds from frozen balance back to
// available balance within a database transaction.
func (r *AccountRepository) UnfreezeBalance(accountID uuid.UUID, amount float64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return apperr.InternalErr("failed to begin transaction", err)
	}
	defer tx.Rollback()

	// Lock the row for update to prevent concurrent balance modifications.
	var frozenBalance float64
	err = tx.QueryRow(
		`SELECT frozen_balance FROM accounts WHERE id = $1 FOR UPDATE`,
		accountID,
	).Scan(&frozenBalance)
	if err != nil {
		return apperr.NotFound("account not found")
	}

	if frozenBalance < amount {
		return apperr.BadRequest("insufficient frozen balance")
	}

	_, err = tx.Exec(
		`UPDATE accounts SET frozen_balance = frozen_balance - $2 WHERE id = $1`,
		accountID, amount,
	)
	if err != nil {
		return apperr.InternalErr("failed to unfreeze balance", err)
	}

	return tx.Commit()
}

// Credit adds funds to an account's available balance.
func (r *AccountRepository) Credit(accountID uuid.UUID, amount float64) error {
	_, err := r.db.Exec(
		`UPDATE accounts SET balance = balance + $2, updated_at = NOW() WHERE id = $1`,
		accountID, amount,
	)
	return err
}

// Debit removes funds from an account's available balance.
func (r *AccountRepository) Debit(accountID uuid.UUID, amount float64) error {
	_, err := r.db.Exec(
		`UPDATE accounts SET balance = balance - $2, updated_at = NOW() WHERE id = $1`,
		accountID, amount,
	)
	return err
}

// RecordBalanceHistory inserts a balance change event into the history table
// and populates the CreatedAt field.
func (r *AccountRepository) RecordBalanceHistory(entry *models.BalanceHistory) error {
	query := `
		INSERT INTO balance_history (id, user_id, account_id, currency, amount, balance_before, balance_after, type, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at`

	return r.db.QueryRow(query,
		entry.ID, entry.UserID, entry.AccountID, entry.Currency,
		entry.Amount, entry.BalanceBefore, entry.BalanceAfter,
		entry.Type, entry.Description,
	).Scan(&entry.CreatedAt)
}
