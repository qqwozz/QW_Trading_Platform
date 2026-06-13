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
	ExpiredAt      *time.Time   `json:"expired_at,omitempty" db:"expired_at"`
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

type MarketTicker struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Symbol       string    `json:"symbol" db:"symbol"`
	LastPrice    float64   `json:"last_price" db:"last_price"`
	BestBid      float64   `json:"best_bid" db:"best_bid"`
	BestAsk      float64   `json:"best_ask" db:"best_ask"`
	BidVolume    float64   `json:"bid_volume" db:"bid_volume"`
	AskVolume    float64   `json:"ask_volume" db:"ask_volume"`
	Volume24h    float64   `json:"volume_24h" db:"volume_24h"`
	High24h      float64   `json:"high_24h" db:"high_24h"`
	Low24h       float64   `json:"low_24h" db:"low_24h"`
	Change24h    float64   `json:"change_24h" db:"change_24h"`
	ChangePct24h float64   `json:"change_pct_24h" db:"change_pct_24h"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type OrderBookLevel struct {
	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity"`
}

type OrderBookSnapshot struct {
	ID        uuid.UUID       `json:"id" db:"id"`
	Symbol    string          `json:"symbol" db:"symbol"`
	Bids      []OrderBookLevel `json:"bids"`
	Asks      []OrderBookLevel `json:"asks"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
}

type TickerResponse struct {
	ID           string  `json:"id"`
	Symbol       string  `json:"symbol"`
	LastPrice    float64 `json:"last_price"`
	BestBid      float64 `json:"best_bid"`
	BestAsk      float64 `json:"best_ask"`
	BidVolume    float64 `json:"bid_volume"`
	AskVolume    float64 `json:"ask_volume"`
	Volume24h    float64 `json:"volume_24h"`
	High24h      float64 `json:"high_24h"`
	Low24h       float64 `json:"low_24h"`
	Change24h    float64 `json:"change_24h"`
	ChangePct24h float64 `json:"change_pct_24h"`
	UpdatedAt    string  `json:"updated_at"`
}

type OrderBookResponse struct {
	Symbol string          `json:"symbol"`
	Bids   []OrderBookLevel `json:"bids"`
	Asks   []OrderBookLevel `json:"asks"`
	Depth  int             `json:"depth"`
}

type BalanceHistoryType string

const (
	BalanceHistoryTypeDeposit    BalanceHistoryType = "DEPOSIT"
	BalanceHistoryTypeWithdrawal BalanceHistoryType = "WITHDRAWAL"
	BalanceHistoryTypeCredit     BalanceHistoryType = "CREDIT"
	BalanceHistoryTypeDebit      BalanceHistoryType = "DEBIT"
	BalanceHistoryTypeFee        BalanceHistoryType = "FEE"
)

type BalanceHistory struct {
	ID           uuid.UUID         `json:"id" db:"id"`
	UserID       uuid.UUID         `json:"user_id" db:"user_id"`
	AccountID    uuid.UUID         `json:"account_id" db:"account_id"`
	Currency     string            `json:"currency" db:"currency"`
	Amount       float64           `json:"amount" db:"amount"`
	BalanceBefore float64          `json:"balance_before" db:"balance_before"`
	BalanceAfter  float64          `json:"balance_after" db:"balance_after"`
	Type         BalanceHistoryType `json:"type" db:"type"`
	Description  string            `json:"description" db:"description"`
	CreatedAt    time.Time         `json:"created_at" db:"created_at"`
}

type PositionHistoryType string

const (
	PositionHistoryTypeOpen     PositionHistoryType = "OPEN"
	PositionHistoryTypeClose    PositionHistoryType = "CLOSE"
	PositionHistoryTypeIncrease PositionHistoryType = "INCREASE"
	PositionHistoryTypeDecrease PositionHistoryType = "DECREASE"
)

type PositionHistory struct {
	ID              uuid.UUID            `json:"id" db:"id"`
	UserID          uuid.UUID            `json:"user_id" db:"user_id"`
	AccountID       uuid.UUID            `json:"account_id" db:"account_id"`
	Symbol          string               `json:"symbol" db:"symbol"`
	QuantityChange  float64              `json:"quantity_change" db:"quantity_change"`
	QuantityBefore  float64              `json:"quantity_before" db:"quantity_before"`
	QuantityAfter   float64              `json:"quantity_after" db:"quantity_after"`
	AvgPriceBefore  float64              `json:"avg_price_before" db:"avg_price_before"`
	AvgPriceAfter   float64              `json:"avg_price_after" db:"avg_price_after"`
	Type            PositionHistoryType  `json:"type" db:"type"`
	TradeID         *uuid.UUID           `json:"trade_id,omitempty" db:"trade_id"`
	CreatedAt       time.Time            `json:"created_at" db:"created_at"`
}
