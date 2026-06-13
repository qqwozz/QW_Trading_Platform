// Package models defines the domain types shared across all services in the
// QW Trading Platform. Types map directly to database tables and JSON APIs.
package models

import (
	"time"

	"github.com/google/uuid"
)

// UserStatus represents the lifecycle state of a user account.
type UserStatus string

const (
	// UserStatusActive indicates a fully operational user account.
	UserStatusActive UserStatus = "ACTIVE"
	// UserStatusSuspended indicates a temporarily disabled user account.
	UserStatusSuspended UserStatus = "SUSPENDED"
	// UserStatusBanned indicates a permanently disabled user account.
	UserStatusBanned UserStatus = "BANNED"
)

// User represents a registered platform user.
type User struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	Email        string     `json:"email" db:"email"`
	Username     string     `json:"username" db:"username"`
	PasswordHash string     `json:"-" db:"password_hash"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	Status       UserStatus `json:"status" db:"status"`
}

// AccountType indicates whether an account is a cash or margin account.
type AccountType string

const (
	// AccountTypeCash is a standard cash account with no leverage.
	AccountTypeCash AccountType = "CASH"
	// AccountTypeMargin is a leveraged margin account.
	AccountTypeMargin AccountType = "MARGIN"
)

// AccountStatus represents the lifecycle state of a trading account.
type AccountStatus string

const (
	// AccountStatusActive indicates a fully operational trading account.
	AccountStatusActive AccountStatus = "ACTIVE"
	// AccountStatusFrozen indicates a temporarily frozen trading account.
	AccountStatusFrozen AccountStatus = "FROZEN"
	// AccountStatusClosed indicates a permanently closed trading account.
	AccountStatusClosed AccountStatus = "CLOSED"
)

// Account represents a trading account that holds a balance in a specific currency.
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

// AssetStatus represents the trading state of an asset.
type AssetStatus string

const (
	// AssetStatusTrading indicates the asset is available for trading.
	AssetStatusTrading AssetStatus = "TRADING"
	// AssetStatusHalted indicates the asset is temporarily suspended from trading.
	AssetStatusHalted AssetStatus = "HALTED"
	// AssetStatusDelisted indicates the asset has been permanently removed.
	AssetStatusDelisted AssetStatus = "DELISTED"
)

// Asset represents a tradeable instrument with its trading parameters.
type Asset struct {
	ID            uuid.UUID   `json:"id" db:"id"`
	Symbol        string      `json:"symbol" db:"symbol"`
	BaseCurrency  string      `json:"base_currency" db:"base_currency"`
	QuoteCurrency string      `json:"quote_currency" db:"quote_currency"`
	MinOrderSize  float64     `json:"min_order_size" db:"min_order_size"`
	TickSize      float64     `json:"tick_size" db:"tick_size"`
	Status        AssetStatus `json:"status" db:"status"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
}

// OrderSide indicates whether an order is to buy or sell.
type OrderSide string

const (
	// OrderSideBuy represents a buy order.
	OrderSideBuy OrderSide = "BUY"
	// OrderSideSell represents a sell order.
	OrderSideSell OrderSide = "SELL"
)

// OrderType indicates whether an order is a limit or market order.
type OrderType string

const (
	// OrderTypeLimit is an order that executes at a specified price or better.
	OrderTypeLimit OrderType = "LIMIT"
	// OrderTypeMarket is an order that executes at the current market price.
	OrderTypeMarket OrderType = "MARKET"
)

// OrderStatus represents the current state of an order in its lifecycle.
type OrderStatus string

const (
	// OrderStatusPending indicates the order has been submitted but not yet active.
	OrderStatusPending OrderStatus = "PENDING"
	// OrderStatusOpen indicates the order is active and waiting to be filled.
	OrderStatusOpen OrderStatus = "OPEN"
	// OrderStatusPartiallyFilled indicates the order has been partially executed.
	OrderStatusPartiallyFilled OrderStatus = "PARTIALLY_FILLED"
	// OrderStatusFilled indicates the order has been fully executed.
	OrderStatusFilled OrderStatus = "FILLED"
	// OrderStatusCancelled indicates the order was cancelled by the user.
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

// TimeInForce specifies how long an order remains active.
type TimeInForce string

const (
	// TimeInForceGTC (Good Till Cancelled) keeps the order active until filled or cancelled.
	TimeInForceGTC TimeInForce = "GTC"
	// TimeInForceIOC (Immediate Or Cancel) fills as much as possible immediately,
	// cancelling any unfilled portion.
	TimeInForceIOC TimeInForce = "IOC"
	// TimeInForceFOK (Fill Or Kill) fills the entire order immediately or cancels it.
	TimeInForceFOK TimeInForce = "FOK"
)

// Order represents a trading order placed by a user.
type Order struct {
	ID             uuid.UUID   `json:"id" db:"id"`
	UserID         uuid.UUID   `json:"user_id" db:"user_id"`
	AccountID      uuid.UUID   `json:"account_id" db:"account_id"`
	Symbol         string      `json:"symbol" db:"symbol"`
	Side           OrderSide   `json:"side" db:"side"`
	Type           OrderType   `json:"type" db:"type"`
	Price          *float64    `json:"price,omitempty" db:"price"`
	Quantity       float64     `json:"quantity" db:"quantity"`
	FilledQuantity float64     `json:"filled_quantity" db:"filled_quantity"`
	Status         OrderStatus `json:"status" db:"status"`
	TimeInForce    TimeInForce `json:"time_in_force" db:"time_in_force"`
	CreatedAt      time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at" db:"updated_at"`
	ExpiredAt      *time.Time  `json:"expired_at,omitempty" db:"expired_at"`
}

// RemainingQuantity returns the unfilled quantity of the order.
func (o *Order) RemainingQuantity() float64 {
	return o.Quantity - o.FilledQuantity
}

// Trade represents an executed trade between a buyer and a seller.
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

// Position represents a user's holdings in a specific asset.
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

// MarketTicker represents real-time market data for a trading pair.
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

// OrderBookLevel represents a single price level in the order book.
type OrderBookLevel struct {
	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity"`
}

// OrderBookSnapshot represents a point-in-time snapshot of the order book for a symbol.
type OrderBookSnapshot struct {
	ID        uuid.UUID        `json:"id" db:"id"`
	Symbol    string           `json:"symbol" db:"symbol"`
	Bids      []OrderBookLevel `json:"bids"`
	Asks      []OrderBookLevel `json:"asks"`
	CreatedAt time.Time        `json:"created_at" db:"created_at"`
}

// TickerResponse is the API response format for market ticker data.
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

// OrderBookResponse is the API response format for order book data.
type OrderBookResponse struct {
	Symbol string           `json:"symbol"`
	Bids   []OrderBookLevel `json:"bids"`
	Asks   []OrderBookLevel `json:"asks"`
	Depth  int              `json:"depth"`
}

// BalanceHistoryType indicates the type of balance change event.
type BalanceHistoryType string

const (
	// BalanceHistoryTypeDeposit records a funds deposit.
	BalanceHistoryTypeDeposit BalanceHistoryType = "DEPOSIT"
	// BalanceHistoryTypeWithdrawal records a funds withdrawal.
	BalanceHistoryTypeWithdrawal BalanceHistoryType = "WITHDRAWAL"
	// BalanceHistoryTypeCredit records a credit (e.g., trade settlement).
	BalanceHistoryTypeCredit BalanceHistoryType = "CREDIT"
	// BalanceHistoryTypeDebit records a debit (e.g., fee deduction).
	BalanceHistoryTypeDebit BalanceHistoryType = "DEBIT"
	// BalanceHistoryTypeFee records a fee charge.
	BalanceHistoryTypeFee BalanceHistoryType = "FEE"
)

// BalanceHistory represents a change in an account's balance over time.
type BalanceHistory struct {
	ID            uuid.UUID          `json:"id" db:"id"`
	UserID        uuid.UUID          `json:"user_id" db:"user_id"`
	AccountID     uuid.UUID          `json:"account_id" db:"account_id"`
	Currency      string             `json:"currency" db:"currency"`
	Amount        float64            `json:"amount" db:"amount"`
	BalanceBefore float64            `json:"balance_before" db:"balance_before"`
	BalanceAfter  float64            `json:"balance_after" db:"balance_after"`
	Type          BalanceHistoryType `json:"type" db:"type"`
	Description   string             `json:"description" db:"description"`
	CreatedAt     time.Time          `json:"created_at" db:"created_at"`
}

// PositionHistoryType indicates the type of position change event.
type PositionHistoryType string

const (
	// PositionHistoryTypeOpen records the opening of a new position.
	PositionHistoryTypeOpen PositionHistoryType = "OPEN"
	// PositionHistoryTypeClose records the closing of a position.
	PositionHistoryTypeClose PositionHistoryType = "CLOSE"
	// PositionHistoryTypeIncrease records an increase in position size.
	PositionHistoryTypeIncrease PositionHistoryType = "INCREASE"
	// PositionHistoryTypeDecrease records a decrease in position size.
	PositionHistoryTypeDecrease PositionHistoryType = "DECREASE"
)

// PositionHistory represents a change in a user's position over time.
type PositionHistory struct {
	ID             uuid.UUID           `json:"id" db:"id"`
	UserID         uuid.UUID           `json:"user_id" db:"user_id"`
	AccountID      uuid.UUID           `json:"account_id" db:"account_id"`
	Symbol         string              `json:"symbol" db:"symbol"`
	QuantityChange float64             `json:"quantity_change" db:"quantity_change"`
	QuantityBefore float64             `json:"quantity_before" db:"quantity_before"`
	QuantityAfter  float64             `json:"quantity_after" db:"quantity_after"`
	AvgPriceBefore float64             `json:"avg_price_before" db:"avg_price_before"`
	AvgPriceAfter  float64             `json:"avg_price_after" db:"avg_price_after"`
	Type           PositionHistoryType `json:"type" db:"type"`
	TradeID        *uuid.UUID          `json:"trade_id,omitempty" db:"trade_id"`
	CreatedAt      time.Time           `json:"created_at" db:"created_at"`
}
