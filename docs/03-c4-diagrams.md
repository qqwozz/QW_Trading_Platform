# C4 Architecture Diagrams

## Context Diagram (Level 1)

```mermaid
C4Context
    title System Context Diagram - QW Trading Platform

    Person(trader, "Trader", "Registered user who trades on the platform")
    Person(admin, "Administrator", "Manages platform operations")

    System(trading_platform, "QW Trading Platform", "Provides virtual trading capabilities with real-time order matching")

    System_Ext(market_data, "External Market Data", "Real-time price feeds from external exchanges")
    System_Ext(monitoring, "Monitoring System", "Prometheus/Grafana for metrics and alerting")

    Rel(trader, trading_platform, "Trades via Web/Mobile")
    Rel(admin, trading_platform, "Manages via Admin Panel")
    Rel(trading_platform, market_data, "Fetches price data", "HTTPS")
    Rel(trading_platform, monitoring, "Exports metrics", "Prometheus")
```

## Container Diagram (Level 2)

```mermaid
C4Container
    title Container Diagram - QW Trading Platform

    Person(trader, "Trader", "Uses web browser to trade")

    System_Boundary(platform, "QW Trading Platform") {
        Container(api_gateway, "API Gateway", "Go", "Routes requests, auth, rate limiting")
        Container(user_service, "User Service", "Go", "User registration, authentication, profile")
        Container(account_service, "Account Service", "Go", "Account management, balance tracking")
        Container(order_service, "Order Service", "Go", "Order lifecycle management")
        Container(matching_engine, "Matching Engine", "C++", "High-performance order matching, order book")
        Container(market_data_service, "Market Data Service", "Go", "Real-time market data via WebSocket")
        Container(analytics_service, "Analytics Service", "Go", "Trading analytics, PnL, metrics")
        ContainerDb(postgres, "PostgreSQL", "Persistent storage for users, accounts, orders, trades")
        ContainerDb(redis, "Redis", "Caching, session store, rate limiting")
        ContainerQueue(kafka, "Kafka", "Event streaming between services")
    }

    Rel(trader, api_gateway, "HTTPS/WSS")
    Rel(api_gateway, user_service, "gRPC")
    Rel(api_gateway, account_service, "gRPC")
    Rel(api_gateway, order_service, "gRPC")
    Rel(api_gateway, market_data_service, "WebSocket")
    Rel(order_service, matching_engine, "gRPC")
    Rel(matching_engine, market_data_service, "Publishes trades", "Kafka")
    Rel(order_service, kafka, "Publishes events")
    Rel(account_service, kafka, "Consumes events")
    Rel(user_service, postgres, "SQL")
    Rel(account_service, postgres, "SQL")
    Rel(order_service, postgres, "SQL")
    Rel(matching_engine, redis, "Order book state")
    Rel(analytics_service, postgres, "SQL")
    Rel(market_data_service, redis, "Market data cache")
```

## Component Diagram - Matching Engine

```mermaid
C4Component
    title Component Diagram - Matching Engine (C++)

    Container_Ext(order_service, "Order Service", "Go", "Submits orders")

    Component(order_queue, "Order Queue", "Concurrent Queue", "Buffers incoming orders")
    Component(order_book, "Order Book", "C++ Core", "Maintains bid/ask levels")
    Component(matcher, "Matcher", "C++ Core", "Price-Time Priority matching")
    Component(fee_calculator, "Fee Calculator", "C++ Core", "Calculates trading fees")
    Component(event_publisher, "Event Publisher", "C++ Core", "Publishes execution events")

    Rel(order_service, order_queue, "New orders", "gRPC")
    Rel(order_queue, order_book, "Enqueue orders")
    Rel(order_book, matcher, "Finds matches")
    Rel(matcher, fee_calculator, "Calculates fees")
    Rel(matcher, event_publisher, "Trade events")
    Rel(event_publisher, order_service, "Execution results", "gRPC")
```

## Component Diagram - API Gateway

```mermaid
C4Component
    title Component Diagram - API Gateway (Go)

    Container_Ext(trader, "Trader", "Web/Mobile client")

    Component(rate_limiter, "Rate Limiter", "Go", "Token bucket rate limiting")
    Component(auth_middleware, "Auth Middleware", "Go", "JWT validation")
    Component(router, "Router", "Go", "Request routing")
    Component(ws_handler, "WebSocket Handler", "Go", "Market data streaming")
    Component(rest_handler, "REST Handler", "Go", "REST API endpoints")

    Component_Ext(user_service, "User Service", "Go")
    Component_Ext(account_service, "Account Service", "Go")
    Component_Ext(order_service, "Order Service", "Go")
    Component_Ext(market_data_service, "Market Data Service", "Go")

    Rel(trader, rate_limiter, "HTTPS/WSS")
    Rel(rate_limiter, auth_middleware, "Filtered requests")
    Rel(auth_middleware, router, "Authenticated requests")
    Rel(router, rest_handler, "REST routes")
    Rel(router, ws_handler, "WebSocket upgrade")
    Rel(rest_handler, user_service, "gRPC")
    Rel(rest_handler, account_service, "gRPC")
    Rel(rest_handler, order_service, "gRPC")
    Rel(ws_handler, market_data_service, "gRPC Stream")
```

## Sequence Diagram - Order Placement Flow

```mermaid
sequenceDiagram
    participant T as Trader
    participant GW as API Gateway
    participant OS as Order Service
    participant AS as Account Service
    participant ME as Matching Engine
    participant MDS as Market Data Service

    T->>GW: POST /api/v1/orders
    GW->>GW: Rate Limit Check
    GW->>GW: JWT Validation
    GW->>OS: CreateOrder (gRPC)

    OS->>AS: FreezeBalance (gRPC)
    AS-->>OS: Balance Frozen

    OS->>ME: SubmitOrder (gRPC)
    ME->>ME: Add to Order Book
    ME->>ME: Attempt Match

    alt Match Found
        ME->>ME: Execute Trade
        ME-->>OS: TradeExecuted
        OS->>AS: UpdateBalances (gRPC)
        AS-->>OS: Balances Updated
        OS-->>GW: Order Filled
        GW-->>T: 200 OK (Order + Trade)
        ME->>MDS: Publish Trade
        MDS-->>T: WebSocket Update
    else No Match
        ME-->>OS: OrderPlaced
        OS-->>GW: Order Open
        GW-->>T: 200 OK (Order)
    end
```

## Sequence Diagram - Market Data Flow

```mermaid
sequenceDiagram
    participant T as Trader
    participant GW as API Gateway
    participant MDS as Market Data Service
    participant ME as Matching Engine
    participant Redis as Redis

    T->>GW: WSS /ws/market
    GW->>MDS: Subscribe (gRPC Stream)

    loop Every Trade/Update
        ME->>MDS: TradeEvent (Kafka)
        MDS->>Redis: Update Cache
        MDS->>MDS: Calculate Best Bid/Ask
        MDS-->>GW: MarketUpdate (gRPC Stream)
        GW-->>T: WebSocket Message
    end
```
