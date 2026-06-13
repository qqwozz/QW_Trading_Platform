package models

import (
	"time"

	"github.com/google/uuid"
)

type UserStatus string

const (
	UserStatusActive    UserStatus = "ACTIVE"
	UserStatusSuspended UserStatus = "SUSPENDED"
	UserStatusBanned    UserStatus = "BANNED"
)

type User struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	Email        string     `json:"email" db:"email"`
	Username     string     `json:"username" db:"username"`
	PasswordHash string     `json:"-" db:"password_hash"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	Status       UserStatus `json:"status" db:"status"`
}

type AccountType string

const (
	AccountTypeCash   AccountType = "CASH"
	AccountTypeMargin AccountType = "MARGIN"
)

type AccountStatus string

const (
	AccountStatusActive AccountStatus = "ACTIVE"
	AccountStatusFrozen AccountStatus = "FROZEN"
	AccountStatusClosed AccountStatus = "CLOSED"
)

type Account struct {
	ID            uuid.UUID     `json:"id" db:"id"`
	UserID        uuid.UUID     `json:"user_id" db:"user_id"`
	Type          AccountType   `json:"type" db:"type"`
	Balance       float64       `json:"balance" db:"balance"`
	FrozenBalance float64       `json:"frozen_balance" db:"frozen_balance"`
	Currency      string        `json:"currency" db:"currency"`
	CreatedAt     time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at" db:"updated_at"`
	Status        AccountStatus `json:"status" db:"status"`
}

type AssetStatus string

const (
	AssetStatusTrading  AssetStatus = "TRADING"
	AssetStatusHalted   AssetStatus = "HALTED"
	AssetStatusDelisted AssetStatus = "DELISTED"
)

type Asset struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	Symbol        string     `json:"symbol" db:"symbol"`
	BaseCurrency  string     `json:"base_currency" db:"base_currency"`
	QuoteCurrency string     `json:"quote_currency" db:"quote_currency"`
	MinOrderSize  float64    `json:"min_order_size" db:"min_order_size"`
	TickSize      float64    `json:"tick_size" db:"tick_size"`
	Status        AssetStatus `json:"status" db:"status"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
}

type OrderSide string

const (
	OrderSideBuy  OrderSide = "BUY"
	OrderSideSell OrderSide = "SELL"
)

type OrderType string

const (
	OrderTypeLimit  OrderType = "LIMIT"
	OrderTypeMarket OrderType = "MARKET"
)

type OrderStatus string

const (
	OrderStatusPending         OrderStatus = "PENDING"
	OrderStatusOpen            OrderStatus = "OPEN"
	OrderStatusPartiallyFilled OrderStatus = "PARTIALLY_FILLED"
	OrderStatusFilled          OrderStatus = "FILLED"
	OrderStatusCancelled       OrderStatus = "CANCELLED"
)

type TimeInForce string

const (
	TimeInForceGTC TimeInForce = "GTC"
	TimeInForceIOC TimeInForce = "IOC"
	TimeInForceFOK TimeInForce = "FOK"
)

type Order struct {
	ID             uuid.UUID    `json:"id" db:"id"`
	UserID         uuid.UUID    `json:"user_id" db:"user_id"`
	AccountID      uuid.UUID    `json:"account_id" db:"account_id"`
	Symbol         string       `json:"symbol" db:"symbol"`
	Side           OrderSide    `json:"side" db:"side"`
	Type           OrderType    `json:"type" db:"type"`
	Price          *float64     `json:"price,omitempty" db:"price"`
	Quantity       float64      `json:"quantity" db:"quantity"`
	FilledQuantity float64      `json:"filled_quantity" db:"filled_quantity"`
	Status         OrderStatus  `json:"status" db:"status"`
	TimeInForce    TimeInForce  `json:"time_in_force" db:"time_in_force"`
	CreatedAt      time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at" db:"updated_at"`
}

func (o *Order) RemainingQuantity() float64 {
	return o.Quantity - o.FilledQuantity
}

type Trade struct {
	ID            uuid.UUID `json:"id" db:"id"`
	Symbol        string    `json:"symbol" db:"symbol"`
	BuyerOrderID  uuid.UUID `json:"buyer_order_id" db:"buyer_order_id"`
	SellerOrderID uuid.UUID `json:"seller_order_id" db:"seller_order_id"`
	BuyerID       uuid.UUID `json:"buyer_id" db:"buyer_id"`
	SellerID      uuid.UUID `json:"seller_id" db:"seller_id"`
	Price         float64   `json:"price" db:"price"`
	Quantity      float64   `json:"quantity" db:"quantity"`
	BuyerFee      float64   `json:"buyer_fee" db:"buyer_fee"`
	SellerFee     float64   `json:"seller_fee" db:"seller_fee"`
	ExecutedAt    time.Time `json:"executed_at" db:"executed_at"`
}

type Position struct {
	ID            uuid.UUID `json:"id" db:"id"`
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	AccountID     uuid.UUID `json:"account_id" db:"account_id"`
	Symbol        string    `json:"symbol" db:"symbol"`
	Quantity      float64   `json:"quantity" db:"quantity"`
	AveragePrice  float64   `json:"average_price" db:"average_price"`
	UnrealizedPnL float64   `json:"unrealized_pnl" db:"unrealized_pnl"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}
