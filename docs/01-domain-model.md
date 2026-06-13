# Domain Model

## Business Entities

### User
Represents a registered platform participant.

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Unique identifier |
| email | string | Email address (unique) |
| username | string | Display name |
| password_hash | string | Bcrypt hash |
| created_at | timestamp | Registration time |
| status | enum | ACTIVE, SUSPENDED, BANNED |

### Account
Virtual trading account linked to a user.

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Unique identifier |
| user_id | UUID | Owner reference |
| type | enum | CASH, MARGIN |
| balance | decimal | Available balance |
| frozen_balance | decimal | Reserved for orders |
| currency | string | ISO 4217 code |
| created_at | timestamp | Creation time |
| status | enum | ACTIVE, FROZEN, CLOSED |

### Asset
Tradeable instrument on the platform.

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Unique identifier |
| symbol | string | Trading pair (e.g., BTC/USDT) |
| base_currency | string | Base asset |
| quote_currency | string | Quote asset |
| min_order_size | decimal | Minimum order quantity |
| tick_size | decimal | Minimum price increment |
| status | enum | TRADING, HALTED, DELISTED |

### Order
A buy or sell instruction.

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Unique identifier |
| user_id | UUID | Owner reference |
| account_id | UUID | Account reference |
| symbol | string | Trading pair |
| side | enum | BUY, SELL |
| type | enum | LIMIT, MARKET |
| price | decimal | Limit price (null for market) |
| quantity | decimal | Original quantity |
| filled_quantity | decimal | Executed quantity |
| status | enum | PENDING, OPEN, PARTIALLY_FILLED, FILLED, CANCELLED |
| time_in_force | enum | GTC, IOC, FOK |
| created_at | timestamp | Order creation time |
| updated_at | timestamp | Last update time |

### Trade
An executed transaction between two orders.

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Unique identifier |
| symbol | string | Trading pair |
| buyer_order_id | UUID | Buy order reference |
| seller_order_id | UUID | Sell order reference |
| buyer_id | UUID | Buyer reference |
| seller_id | UUID | Seller reference |
| price | decimal | Execution price |
| quantity | decimal | Executed quantity |
| buyer_fee | decimal | Buyer commission |
| seller_fee | decimal | Seller commission |
| executed_at | timestamp | Execution time |

### Position
Tracks holdings of a specific asset.

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Unique identifier |
| user_id | UUID | Owner reference |
| account_id | UUID | Account reference |
| symbol | string | Asset symbol |
| quantity | decimal | Current quantity |
| average_price | decimal | Weighted average entry price |
| unrealized_pnl | decimal | Current P&L |
| updated_at | timestamp | Last update time |

## Value Objects

### Price
- Precision: 18 digits, 8 decimal places
- Validation: must be positive for orders

### Quantity
- Precision: 18 digits, 8 decimal places
- Validation: must be positive, >= min_order_size

### Money
- Amount + Currency pair
- Used for balance operations

## Aggregates

1. **User Aggregate**: User + Accounts
2. **Order Aggregate**: Order lifecycle management
3. **Trade Aggregate**: Trade execution records
4. **Position Aggregate**: Position + balance updates

## Domain Events

| Event | Aggregate | Description |
|-------|-----------|-------------|
| UserRegistered | User | New user created |
| AccountCreated | User | New account opened |
| BalanceDeposited | User | Funds added |
| BalanceWithdrawn | User | Funds removed |
| OrderCreated | Order | New order submitted |
| OrderCancelled | Order | Order cancelled |
| OrderFilled | Order | Order fully executed |
| OrderPartiallyFilled | Order | Order partially executed |
| TradeExecuted | Trade | Trade matched and executed |
| PositionOpened | Position | New position created |
| PositionUpdated | Position | Position quantity changed |
| PositionClosed | Position | Position zeroed out |

## Business Rules

1. A user must have at least one account to trade
2. Orders can only be placed with sufficient available balance
3. SELL orders require sufficient asset quantity
4. Price-Time Priority determines match order
5. Partial fills are allowed for LIMIT orders
6. MARKET orders execute at best available price
7. Cancelled orders release frozen funds
8. Positions are calculated from trade history
