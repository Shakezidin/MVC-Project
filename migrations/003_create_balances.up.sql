CREATE TABLE IF NOT EXISTS balances (
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id        UUID NOT NULL UNIQUE REFERENCES bank_accounts(id) ON DELETE CASCADE,
    available_balance DECIMAL(18, 2) NOT NULL DEFAULT 0.00,
    current_balance   DECIMAL(18, 2) NOT NULL DEFAULT 0.00,
    currency          VARCHAR(10) NOT NULL DEFAULT 'INR',
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_balances_account_id ON balances (account_id);
