# ER Diagram

```mermaid
erDiagram
    USER {
        uuid id PK
        string email UK
        string username
        string password_hash
        timestamp created_at
        string status
    }

    ACCOUNT {
        uuid id PK
        uuid user_id FK
        string type
        decimal balance
        decimal frozen_balance
        string currency
        timestamp created_at
        string status
    }

    ASSET {
        uuid id PK
        string symbol UK
        string base_currency
        string quote_currency
        decimal min_order_size
        decimal tick_size
        string status
    }

    ORDER {
        uuid id PK
        uuid user_id FK
        uuid account_id FK
        string symbol FK
        string side
        string type
        decimal price
        decimal quantity
        decimal filled_quantity
        string status
        string time_in_force
        timestamp created_at
        timestamp updated_at
    }

    TRADE {
        uuid id PK
        string symbol FK
        uuid buyer_order_id FK
        uuid seller_order_id FK
        uuid buyer_id FK
        uuid seller_id FK
        decimal price
        decimal quantity
        decimal buyer_fee
        decimal seller_fee
        timestamp executed_at
    }

    POSITION {
        uuid id PK
        uuid user_id FK
        uuid account_id FK
        string symbol FK
        decimal quantity
        decimal average_price
        decimal unrealized_pnl
        timestamp updated_at
    }

    ORDER_BOOK {
        uuid id PK
        string symbol FK
        string side
        decimal price
        decimal total_quantity
        timestamp updated_at
    }

    ORDER_BOOK_ENTRY {
        uuid id PK
        uuid order_book_id FK
        uuid order_id FK
        decimal price
        decimal quantity
        int priority
        timestamp created_at
    }

    USER ||--o{ ACCOUNT : has
    USER ||--o{ ORDER : places
    USER ||--o{ POSITION : holds
    ACCOUNT ||--o{ ORDER : funds
    ACCOUNT ||--o{ POSITION : tracks
    ASSET ||--o{ ORDER : trades
    ASSET ||--o{ POSITION : represents
    ASSET ||--o{ ORDER_BOOK : maintains
    ORDER ||--o{ TRADE : "buy side"
    ORDER ||--o{ TRADE : "sell side"
    ORDER_BOOK ||--o{ ORDER_BOOK_ENTRY : contains
    ORDER ||--o{ ORDER_BOOK_ENTRY : "linked to"
```

## Indexes

### USER
- `idx_user_email` UNIQUE on email
- `idx_user_username` on username

### ACCOUNT
- `idx_account_user_id` on user_id
- `idx_account_user_currency` on (user_id, currency) UNIQUE

### ORDER
- `idx_order_user_id` on user_id
- `idx_order_symbol_status` on (symbol, status)
- `idx_order_account_id` on account_id
- `idx_order_created_at` on created_at

### TRADE
- `idx_trade_symbol` on symbol
- `idx_trade_executed_at` on executed_at
- `idx_trade_buyer` on buyer_id
- `idx_trade_seller` on seller_id

### POSITION
- `idx_position_user_symbol` on (user_id, symbol) UNIQUE
- `idx_position_account` on account_id

### ORDER_BOOK
- `idx_orderbook_symbol_side` on (symbol, side)

### ORDER_BOOK_ENTRY
- `idx_obentry_book_id` on order_book_id
- `idx_obentry_order_id` on order_id
