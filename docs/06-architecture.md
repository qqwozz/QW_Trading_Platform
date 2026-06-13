# Architecture Documentation

## Overview

QW Trading Platform is a virtual cryptocurrency exchange simulator built with microservice architecture. The system supports real-time order matching, portfolio management, and market data streaming.

## Technology Stack

| Component | Technology | Rationale |
|-----------|-----------|-----------|
| API Gateway | Go | High concurrency, low latency |
| User Service | Go | Standard CRUD operations |
| Account Service | Go | Financial operations with ACID |
| Order Service | Go | Order lifecycle management |
| Matching Engine | C++ | Maximum performance for order matching |
| Market Data Service | Go | WebSocket handling, real-time streaming |
| Analytics Service | Go | Batch processing, calculations |
| Database | PostgreSQL | ACID compliance, financial data |
| Cache | Redis | Fast access, session management |
| Message Queue | Kafka | Event streaming, decoupling |
| API Protocol | gRPC | High performance, type-safe |
| External API | REST | Client-facing HTTP API |
| Streaming | WebSocket | Real-time market data |

## Service Responsibilities

### API Gateway
- Rate limiting (token bucket)
- JWT authentication
- Request routing
- WebSocket upgrade
- API versioning

### User Service
- User registration
- Authentication (JWT)
- Profile management
- Token refresh

### Account Service
- Account creation
- Balance management
- Fund freezing/unfreezing
- Position tracking

### Order Service
- Order validation
- Order lifecycle (create, cancel, fill)
- Order history
- Trade recording

### Matching Engine (C++)
- Order book management
- Price-Time Priority matching
- Partial/full fill execution
- Fee calculation
- Low-latency processing

### Market Data Service
- Real-time price updates
- Order book snapshots
- Trade streaming
- WebSocket connections

### Analytics Service
- PnL calculation
- ROI computation
- Sharpe ratio
- Max drawdown
- Win rate

## Data Flow

1. **Order Placement**: Trader → Gateway → Order Service → Account Service (freeze) → Matching Engine
2. **Trade Execution**: Matching Engine → Order Service → Account Service (update) → Market Data Service → Trader
3. **Market Data**: Matching Engine → Kafka → Market Data Service → WebSocket → Trader

## Consistency Model

| Component | Consistency | Reason |
|-----------|------------|--------|
| Matching Engine | Strong | Financial accuracy critical |
| Account Service | Strong | Balance integrity |
| Order Service | Strong | Order state management |
| Analytics Service | Eventual | Can tolerate delays |
| Market Data | Eventual | Price can lag slightly |

## Failure Handling

1. **Service Crash**: Kubernetes restarts, state rebuilt from events
2. **Database Failure**: Automatic failover to replica
3. **Kafka Failure**: Local buffer, retry with backoff
4. **Network Partition**: Circuit breaker pattern, graceful degradation

## Security

1. **Authentication**: JWT with refresh tokens
2. **Authorization**: User can only access own data
3. **Encryption**: TLS for all communications
4. **Secrets**: Kubernetes Secrets, Vault integration
5. **Rate Limiting**: Token bucket per IP/user

## Performance Targets

| Metric | Target |
|--------|--------|
| Order placement latency | < 100ms |
| Market data latency | < 50ms |
| WebSocket update latency | < 100ms |
| Orders per second | 1000+ |
| Concurrent users | 10,000+ |
| Uptime | 99.9% |
