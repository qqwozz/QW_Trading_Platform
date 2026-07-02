<div align="center">

# QW Trading Platform

**Production-grade virtual cryptocurrency exchange built on microservices architecture**

[![CI](https://github.com/qw-trading/platform/actions/workflows/ci.yml/badge.svg)](https://github.com/qw-trading/platform/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go)](https://go.dev/)
[![C++](https://img.shields.io/badge/C++-17-00599C?logo=cplusplus)](https://isocpp.org/)
[![React](https://img.shields.io/badge/React-19-61DAFB?logo=react)](https://react.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

<br>

A full-stack trading platform featuring a high-performance C++ matching engine, Go microservices, real-time WebSocket data, and a React trading terminal — all containerized with Docker.

[Architecture](#architecture) · [Quick Start](#quick-start) · [API Reference](#api-reference) · [Documentation](#documentation)

</div>

---

## Features

| Category | Details |
|---|---|
| **Matching Engine** | C++17 in-memory order book with price-time priority, LIMIT/MARKET orders, GTC/IOC/FOK policies, partial fills, and configurable maker/taker fees |
| **Microservices** | 7 independent Go services — API Gateway, User, Account, Order, Portfolio, Market Data, History |
| **Real-time Data** | WebSocket streaming for order book updates and market tickers via gorilla/websocket |
| **Security** | JWT auth (access + refresh tokens), bcrypt password hashing, brute-force protection, per-user data isolation, rate limiting |
| **Trading Terminal** | React 19 SPA with dark theme, live order book, portfolio tracking, and trade history |
| **Infrastructure** | Docker Compose with 9 containers, GitHub Actions CI (lint → test → build → Docker image), PostgreSQL 16, Redis 7 |
| **Pre-configured Pairs** | BTC/USDT, ETH/USDT, SOL/USDT |

---

## Architecture

```
                           ┌──────────────────────────────────────┐
                           │           Frontend (React)           │
                           │          localhost:3000               │
                           └──────────────┬───────────────────────┘
                                          │ /v1/*
                           ┌──────────────▼───────────────────────┐
                           │         API Gateway (:8080)          │
                           │   CORS · Rate Limit · Request ID     │
                           └──┬─────┬──────┬──────┬──────┬───────┘
                              │     │      │      │      │
                 ┌────────────┘     │      │      │      └────────────┐
                 │                  │      │      │                   │
          ┌──────▼──────┐  ┌───────▼────┐ │ ┌────▼─────┐    ┌───────▼───────┐
          │ User Service│  │  Account   │ │ │ Portfolio│    │ History Service│
          │   (:8081)   │  │ Service    │ │ │ Service  │    │    (:8086)     │
          │ Register    │  │  (:8082)   │ │ │ (:8084)  │    │ Orders/Trades  │
          │ Login (JWT) │  │ Deposits   │ │ │ Positions│    │ Balance/Pos    │
          └──────┬──────┘  └───────┬────┘ │ └────┬─────┘    └───────┬───────┘
                 │                 │      │      │                   │
                 │          ┌──────▼──────▼──────▼──────┐           │
                 │          │      Order Service        │           │
                 │          │         (:8083)           │           │
                 │          │   CRUD · Validation       │           │
                 │          └──────────┬────────────────┘           │
                 │                     │                            │
                 │          ┌──────────▼────────────────┐           │
                 │          │   Matching Engine (C++)   │           │
                 │          │         (:50051)          │           │
                 │          │  In-memory order book     │           │
                 │          │  Price-time priority      │           │
                 │          │  Maker/Taker fees         │           │
                 │          └──────────┬────────────────┘           │
                 │                     │                            │
                 │          ┌──────────▼────────────────┐           │
                 │          │  Market Data Service      │           │
                 │          │         (:8085)           │           │
                 │          │  Tickers · Order Book     │           │
                 │          │  WebSocket Hub            │           │
                 │          └──────────┬────────────────┘           │
                 │                     │                            │
                 └─────────────────────┼────────────────────────────┘
                                       │
                              ┌────────▼────────┐
                              │   PostgreSQL 16  │
                              │   Redis 7        │
                              └──────────────────┘
```

---

## Tech Stack

### Backend

| Component | Technology | Purpose |
|---|---|---|
| Microservices | Go 1.25 | Business logic, REST API, WebSocket |
| Matching Engine | C++17 | High-performance order matching |
| Database | PostgreSQL 16 | Persistent storage, ACID transactions |
| Cache | Redis 7 | Rate limiting, session caching |
| Auth | JWT (`golang-jwt/jwt/v5`) | Stateless authentication |
| Password Hashing | bcrypt (`golang.org/x/crypto`) | Secure credential storage |
| DB Driver | `lib/pq` (raw `database/sql`) | Zero-ORM database access |
| WebSocket | `gorilla/websocket` | Real-time market data streaming |
| UUID | `google/uuid` | Unique entity identifiers |

### Frontend

| Component | Technology |
|---|---|
| Framework | React 19 + TypeScript 6 |
| Build Tool | Vite 8 |
| Styling | Tailwind CSS 4 |
| Charts | Recharts 3 |
| Icons | Lucide React |
| Routing | React Router 7 |
| Linting | ESLint 10 + typescript-eslint |

### Infrastructure

| Component | Technology |
|---|---|
| Containerization | Docker + Docker Compose |
| CI/CD | GitHub Actions (lint → test → build → Docker) |
| Linting | `go vet` + `staticcheck` |
| Testing | `go test -race -coverprofile` |

---

## Quick Start

### Docker (Recommended)

```bash
git clone https://github.com/qw-trading/platform.git
cd platform
docker-compose up --build
```

The platform will be available at:

| Service | URL |
|---|---|
| Trading Terminal | http://localhost:3000 |
| API Gateway | http://localhost:8080 |
| Matching Engine | localhost:50051 |

### Local Development

**Prerequisites:** Go 1.25+, Node.js 20+, PostgreSQL 16+, CMake 3.20+

```bash
# 1. Start PostgreSQL
docker run -d --name postgres -p 5432:5432 \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=qw_trading \
  postgres:16-alpine

# 2. Apply migrations
psql -h localhost -U postgres -d qw_trading \
  -f deploy/migrations/001_init.sql \
  -f deploy/migrations/002_audit.sql \
  -f deploy/migrations/003_market_data.sql

# 3. Set environment variables
export DB_HOST=localhost
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=qw_trading
export JWT_SECRET=your-dev-secret

# 4. Start backend services (each in a separate terminal)
go run cmd/api-gateway/main.go
go run cmd/user-service/main.go
go run cmd/account-service/main.go
go run cmd/order-service/main.go
go run cmd/portfolio-service/main.go
go run cmd/market-data-service/main.go
go run cmd/history-service/main.go

# 5. Build the C++ matching engine
mkdir -p build && cd build
cmake .. && make
./matching-engine

# 6. Start the frontend
cd frontend
npm install
npm run dev
```

### Environment Variables

| Variable | Default | Description |
|---|---|---|
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `postgres` | Database user |
| `DB_PASSWORD` | `postgres` | Database password |
| `DB_NAME` | `qw_trading` | Database name |
| `JWT_SECRET` | `change-me-in-production` | JWT signing secret |
| `PORT` | `8080` | Service port |
| `APP_ENV` | `development` | Environment (development/production) |
| `ALLOWED_ORIGINS` | `*` | CORS allowed origins |
| `RATE_LIMIT_RPS` | `100` | Rate limit requests per second |
| `RATE_LIMIT_BURST` | `200` | Rate limit burst capacity |
| `MAKER_FEE_BPS` | `10` | Maker fee in basis points (0.10%) |
| `TAKER_FEE_BPS` | `20` | Taker fee in basis points (0.20%) |
| `JWT_EXPIRY_HOURS` | `1` | Access token TTL |
| `REFRESH_EXPIRY_HOURS` | `168` | Refresh token TTL (7 days) |

---

## API Reference

All endpoints are prefixed with `/v1` and routed through the API Gateway at port **8080**.

### Public Endpoints

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/v1/auth/register` | Create a new account |
| `POST` | `/v1/auth/login` | Authenticate and receive tokens |
| `GET` | `/v1/market/tickers` | All trading pair tickers |
| `GET` | `/v1/market/tickers/{symbol}` | Single pair ticker |
| `GET` | `/v1/market/orderbook/{symbol}` | Order book snapshot |
| `GET` | `/v1/market/ws` | WebSocket for real-time market data |

### Protected Endpoints (JWT Required)

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/v1/users/me` | Current user profile |
| `GET` | `/v1/accounts` | List user accounts |
| `POST` | `/v1/accounts/deposit` | Deposit virtual funds |
| `POST` | `/v1/accounts/withdraw` | Withdraw virtual funds |
| `POST` | `/v1/orders` | Place a new order |
| `GET` | `/v1/orders` | List orders (filterable) |
| `GET` | `/v1/orders/{id}` | Get order details |
| `DELETE` | `/v1/orders/{id}` | Cancel an open order |
| `GET` | `/v1/portfolio` | Portfolio summary with equity |
| `GET` | `/v1/positions` | List open positions |
| `POST` | `/v1/positions` | Manually update a position |
| `GET` | `/v1/balances` | Account balances with frozen breakdown |
| `GET` | `/v1/history/orders` | Order history (paginated) |
| `GET` | `/v1/history/trades` | Trade history (paginated) |
| `GET` | `/v1/history/balance` | Balance change history |
| `GET` | `/v1/history/positions` | Position change history |

### Health Check

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/health` | Service status + DB connectivity |

### Quick Examples

```bash
# Register
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"trader@example.com","username":"trader","password":"secret1234"}'

# Login
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"trader@example.com","password":"secret1234"}'

# Deposit funds
curl -X POST http://localhost:8080/v1/accounts/deposit \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"currency":"USDT","amount":10000}'

# Place a limit order
curl -X POST http://localhost:8080/v1/orders \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"symbol":"BTC/USDT","side":"BUY","type":"LIMIT","price":50000,"quantity":0.1}'

# Get portfolio
curl http://localhost:8080/v1/portfolio \
  -H "Authorization: Bearer <token>"
```

Full OpenAPI specification: [`docs/openapi.yaml`](docs/openapi.yaml)

---

## Project Structure

```
qw-trading-platform/
├── cmd/                            # Service entrypoints
│   ├── api-gateway/
│   ├── user-service/
│   ├── account-service/
│   ├── order-service/
│   ├── portfolio-service/
│   ├── market-data-service/
│   ├── history-service/
│   └── matching-engine/
├── internal/                       # Business logic
│   ├── models/                     # Domain models
│   ├── db/                         # Database layer
│   ├── user/                       # User service
│   ├── account/                    # Account service
│   ├── order/                      # Order service
│   ├── portfolio/                  # Portfolio service
│   ├── market/                     # Market data service
│   └── history/                    # History service
├── pkg/                            # Shared packages
│   ├── config/                     # Environment configuration
│   ├── errors/                     # Error types
│   ├── middleware/                  # CORS, rate limit, auth, logging
│   └── response/                   # HTTP response helpers
├── matching-engine/                # C++ matching engine
│   ├── src/                        # Source files
│   ├── include/                    # Headers
│   └── CMakeLists.txt              # CMake build config
├── frontend/                       # React trading terminal
│   ├── src/
│   │   ├── components/             # UI components
│   │   ├── pages/                  # Route pages
│   │   ├── lib/                    # API client, utilities
│   │   └── App.tsx
│   ├── package.json
│   └── vite.config.ts
├── deploy/                         # Deployment configs
│   ├── Dockerfile-go
│   ├── Dockerfile-cpp
│   ├── Dockerfile-frontend
│   └── migrations/                 # SQL migrations
│       ├── 001_init.sql
│       ├── 002_audit.sql
│       └── 003_market_data.sql
├── docs/                           # Documentation
│   ├── architecture.md
│   ├── development.md
│   ├── openapi.yaml
│   └── ...
├── .github/workflows/ci.yml        # CI pipeline
├── docker-compose.yml
├── go.mod
└── go.sum
```

---

## Database Schema

### Core Tables (`001_init.sql`)

| Table | Description |
|---|---|
| `users` | User accounts (UUID PK, email, username, bcrypt hash, status) |
| `accounts` | Currency accounts per user (balance, frozen_balance, type) |
| `assets` | Trading pair definitions (symbol, min_order_size, tick_size) |
| `orders` | All orders (side, type, price, quantity, filled_quantity, status) |
| `trades` | Executed trades (price, quantity, fees for buyer/seller) |
| `positions` | User positions per symbol (quantity, avg_price, unrealized_pnl) |

### Audit Tables (`002_audit.sql`)

| Table | Description |
|---|---|
| `balance_history` | Every balance change with before/after amounts |
| `position_history` | Every position change with before/after values |

### Market Data Tables (`003_market_data.sql`)

| Table | Description |
|---|---|
| `market_tickers` | Real-time ticker data per symbol |
| `order_book_snapshots` | JSONB order book snapshots |

All financial fields use `NUMERIC(20,8)` for precision.

---

## Services

| Service | Port | Responsibility |
|---|---|---|
| **API Gateway** | `:8080` | HTTP reverse proxy, CORS, rate limiting, request ID injection |
| **User Service** | `:8081` | Registration, login, JWT generation, profile management |
| **Account Service** | `:8082` | Virtual currency accounts, deposits, withdrawals, balance queries |
| **Order Service** | `:8083` | Order CRUD, validation, cancellation |
| **Portfolio Service** | `:8084` | Position tracking, portfolio summary, equity calculations |
| **Market Data Service** | `:8085` | Ticker data, order book snapshots, WebSocket hub for real-time streaming |
| **History Service** | `:8086` | Audit trail — order, trade, balance, and position history |
| **Matching Engine** | `:50051` | In-memory order book, price-time priority matching, fee calculation |
| **Frontend** | `:3000` | React SPA trading terminal served via Nginx |

---

## Matching Engine (C++)

The core matching engine is written in C++17 for maximum performance:

- **In-memory order book** per trading symbol
- **Price-time priority** — orders at the same price matched in FIFO order
- **Order types:** LIMIT, MARKET
- **Time-in-force:** GTC (Good Till Cancelled), IOC (Immediate or Cancel), FOK (Fill or Kill)
- **Partial fills** supported — unfilled quantity remains in the book
- **Fee model:** Maker 0.10% (10 bps) / Taker 0.20% (20 bps), configurable via environment
- **Thread safety:** Per-order-book mutex + engine-level registry mutex
- **Event callbacks:** `on_trade` and `on_order` fire after each matching operation

---

## Testing

```bash
# Run all tests
go test ./...

# With race detection and coverage
go test -v -race -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out
```

Test coverage includes:
- Handler tests with mock repositories for all services
- Unit tests for middleware (rate limiting, CORS, auth, logging)
- Error type and response helper tests
- Configuration and fee calculation tests

---

## CI/CD Pipeline

GitHub Actions workflow (`.github/workflows/ci.yml`) runs on every push/PR to `main`:

```
┌─────────┐     ┌─────────┐     ┌─────────┐     ┌─────────────┐
│  Lint   │────▶│  Test   │────▶│  Build  │────▶│ Docker Build│
│go vet   │     │go test  │     │go build │     │ image push  │
│staticchk│     │-race    │     │7 services│    │             │
└─────────┘     └─────────┘     └─────────┘     └─────────────┘
```

---

## Documentation

| Document | Description |
|---|---|
| [`docs/architecture.md`](docs/architecture.md) | System architecture overview |
| [`docs/development.md`](docs/development.md) | Developer quick-start guide |
| [`docs/openapi.yaml`](docs/openapi.yaml) | OpenAPI 3.0.3 specification |
| [`docs/01-domain-model.md`](docs/01-domain-model.md) | Domain entities and business rules |
| [`docs/02-er-diagram.md`](docs/02-er-diagram.md) | Mermaid ER diagram |
| [`docs/03-c4-diagrams.md`](docs/03-c4-diagrams.md) | C4 Context, Container, and Sequence diagrams |
| [`docs/04-api-contracts.md`](docs/04-api-contracts.md) | REST and WebSocket API contracts |
| [`docs/05-event-flow.md`](docs/05-event-flow.md) | Event-driven architecture design |

---

## License

[MIT](LICENSE) — Dima Kiselev
