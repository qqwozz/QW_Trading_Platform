CREATE TABLE balance_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    currency VARCHAR(10) NOT NULL,
    amount NUMERIC(20,8) NOT NULL,
    balance_before NUMERIC(20,8) NOT NULL,
    balance_after NUMERIC(20,8) NOT NULL,
    type VARCHAR(20) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_balance_history_user_id ON balance_history(user_id);
CREATE INDEX idx_balance_history_account_id ON balance_history(account_id);
CREATE INDEX idx_balance_history_created_at ON balance_history(created_at DESC);

CREATE TABLE position_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    symbol VARCHAR(20) NOT NULL,
    quantity_change NUMERIC(20,8) NOT NULL,
    quantity_before NUMERIC(20,8) NOT NULL,
    quantity_after NUMERIC(20,8) NOT NULL,
    avg_price_before NUMERIC(20,8) NOT NULL,
    avg_price_after NUMERIC(20,8) NOT NULL,
    type VARCHAR(20) NOT NULL,
    trade_id UUID,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_position_history_user_id ON position_history(user_id);
CREATE INDEX idx_position_history_symbol ON position_history(symbol);
CREATE INDEX idx_position_history_created_at ON position_history(created_at DESC);

ALTER TABLE orders ADD COLUMN expired_at TIMESTAMPTZ;
