# API Contracts

## Base URL
```
https://api.qw-trading.local/v1
```

## Authentication
All protected endpoints require:
```
Authorization: Bearer <jwt_access_token>
```

## Rate Limits
- Public endpoints: 100 requests/minute
- Authenticated endpoints: 1000 requests/minute
- Order placement: 100 requests/minute

---

## User Service

### POST /auth/register
Register a new user.

**Request:**
```json
{
  "email": "string",
  "username": "string",
  "password": "string"
}
```

**Response 201:**
```json
{
  "id": "uuid",
  "email": "string",
  "username": "string",
  "created_at": "timestamp"
}
```

**Errors:**
- 400: Invalid input
- 409: Email already exists

---

### POST /auth/login
Authenticate user.

**Request:**
```json
{
  "email": "string",
  "password": "string"
}
```

**Response 200:**
```json
{
  "access_token": "string",
  "refresh_token": "string",
  "expires_in": 3600
}
```

**Errors:**
- 401: Invalid credentials

---

### POST /auth/refresh
Refresh access token.

**Request:**
```json
{
  "refresh_token": "string"
}
```

**Response 200:**
```json
{
  "access_token": "string",
  "expires_in": 3600
}
```

---

### GET /users/me
Get current user profile.

**Response 200:**
```json
{
  "id": "uuid",
  "email": "string",
  "username": "string",
  "created_at": "timestamp"
}
```

---

## Account Service

### GET /accounts
List user accounts.

**Response 200:**
```json
{
  "accounts": [
    {
      "id": "uuid",
      "type": "CASH",
      "balance": "10000.00",
      "frozen_balance": "500.00",
      "currency": "USDT",
      "status": "ACTIVE"
    }
  ]
}
```

---

### GET /accounts/{id}
Get account details.

**Response 200:**
```json
{
  "id": "uuid",
  "type": "CASH",
  "balance": "10000.00",
  "frozen_balance": "500.00",
  "currency": "USDT",
  "status": "ACTIVE",
  "created_at": "timestamp"
}
```

---

### POST /accounts/deposit
Deposit funds.

**Request:**
```json
{
  "currency": "USDT",
  "amount": "1000.00"
}
```

**Response 200:**
```json
{
  "account_id": "uuid",
  "new_balance": "11000.00"
}
```

---

## Order Service

### POST /orders
Place a new order.

**Request:**
```json
{
  "symbol": "BTC/USDT",
  "side": "BUY",
  "type": "LIMIT",
  "price": "50000.00",
  "quantity": "0.1",
  "time_in_force": "GTC"
}
```

**Response 201:**
```json
{
  "id": "uuid",
  "symbol": "BTC/USDT",
  "side": "BUY",
  "type": "LIMIT",
  "price": "50000.00",
  "quantity": "0.1",
  "filled_quantity": "0",
  "status": "OPEN",
  "created_at": "timestamp"
}
```

**Errors:**
- 400: Invalid order parameters
- 403: Insufficient balance
- 404: Symbol not found

---

### DELETE /orders/{id}
Cancel an order.

**Response 200:**
```json
{
  "id": "uuid",
  "status": "CANCELLED",
  "cancelled_at": "timestamp"
}
```

**Errors:**
- 404: Order not found
- 409: Order already filled/cancelled

---

### GET /orders
List user orders.

**Query Parameters:**
- `symbol` (optional): Filter by symbol
- `status` (optional): Filter by status
- `limit` (default: 50): Results per page
- `offset` (default: 0): Pagination offset

**Response 200:**
```json
{
  "orders": [
    {
      "id": "uuid",
      "symbol": "BTC/USDT",
      "side": "BUY",
      "type": "LIMIT",
      "price": "50000.00",
      "quantity": "0.1",
      "filled_quantity": "0.05",
      "status": "PARTIALLY_FILLED",
      "created_at": "timestamp"
    }
  ],
  "total": 100,
  "limit": 50,
  "offset": 0
}
```

---

### GET /orders/{id}
Get order details.

**Response 200:**
```json
{
  "id": "uuid",
  "symbol": "BTC/USDT",
  "side": "BUY",
  "type": "LIMIT",
  "price": "50000.00",
  "quantity": "0.1",
  "filled_quantity": "0.05",
  "status": "PARTIALLY_FILLED",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

---

## Market Data Service

### GET /market/symbols
List available trading pairs.

**Response 200:**
```json
{
  "symbols": [
    {
      "symbol": "BTC/USDT",
      "base_currency": "BTC",
      "quote_currency": "USDT",
      "min_order_size": "0.001",
      "tick_size": "0.01",
      "status": "TRADING"
    }
  ]
}
```

---

### GET /market/ticker/{symbol}
Get ticker data.

**Response 200:**
```json
{
  "symbol": "BTC/USDT",
  "last_price": "50000.00",
  "best_bid": "49999.00",
  "best_ask": "50001.00",
  "volume_24h": "1234.56",
  "high_24h": "51000.00",
  "low_24h": "49000.00",
  "change_24h": "2.5"
}
```

---

### GET /market/orderbook/{symbol}
Get order book snapshot.

**Query Parameters:**
- `limit` (default: 20): Depth levels

**Response 200:**
```json
{
  "symbol": "BTC/USDT",
  "bids": [
    ["49999.00", "1.5"],
    ["49998.00", "2.0"]
  ],
  "asks": [
    ["50001.00", "1.2"],
    ["50002.00", "0.8"]
  ]
}
```

---

### GET /market/trades/{symbol}
Get recent trades.

**Query Parameters:**
- `limit` (default: 100): Number of trades

**Response 200:**
```json
{
  "trades": [
    {
      "id": "uuid",
      "price": "50000.00",
      "quantity": "0.1",
      "side": "BUY",
      "executed_at": "timestamp"
    }
  ]
}
```

---

## WebSocket API

### Connection
```
wss://api.qw-trading.local/ws?token=<jwt_token>
```

### Subscribe to Market Data
```json
{
  "type": "subscribe",
  "channel": "ticker",
  "symbol": "BTC/USDT"
}
```

### Ticker Update
```json
{
  "type": "ticker",
  "symbol": "BTC/USDT",
  "last_price": "50000.00",
  "best_bid": "49999.00",
  "best_ask": "50001.00",
  "volume_24h": "1234.56"
}
```

### Order Book Update
```json
{
  "type": "orderbook",
  "symbol": "BTC/USDT",
  "bids": [["49999.00", "1.5"]],
  "asks": [["50001.00", "1.2"]]
}
```

### Trade Update
```json
{
  "type": "trade",
  "symbol": "BTC/USDT",
  "price": "50000.00",
  "quantity": "0.1",
  "side": "BUY",
  "executed_at": "timestamp"
}
```

### Order Update (Private)
```json
{
  "type": "order_update",
  "order_id": "uuid",
  "status": "FILLED",
  "filled_quantity": "0.1",
  "avg_price": "50000.00"
}
```
