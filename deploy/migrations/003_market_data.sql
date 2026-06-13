CREATE TABLE market_tickers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    symbol VARCHAR(20) UNIQUE NOT NULL,
    last_price NUMERIC(20,8) DEFAULT 0,
    best_bid NUMERIC(20,8) DEFAULT 0,
    best_ask NUMERIC(20,8) DEFAULT 0,
    bid_volume NUMERIC(20,8) DEFAULT 0,
    ask_volume NUMERIC(20,8) DEFAULT 0,
    volume_24h NUMERIC(20,8) DEFAULT 0,
    high_24h NUMERIC(20,8) DEFAULT 0,
    low_24h NUMERIC(20,8) DEFAULT 0,
    change_24h NUMERIC(20,8) DEFAULT 0,
    change_pct_24h NUMERIC(10,4) DEFAULT 0,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_market_tickers_symbol ON market_tickers(symbol);

INSERT INTO market_tickers (symbol) VALUES
    ('BTC/USDT'),
    ('ETH/USDT'),
    ('SOL/USDT')
ON CONFLICT (symbol) DO NOTHING;

CREATE TABLE order_book_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    symbol VARCHAR(20) NOT NULL,
    bids_json JSONB NOT NULL DEFAULT '[]'::jsonb,
    asks_json JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_order_book_snapshots_symbol_time ON order_book_snapshots(symbol, created_at DESC);
