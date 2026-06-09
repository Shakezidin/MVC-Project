CREATE TABLE IF NOT EXISTS beneficiaries (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id          UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    beneficiary_name VARCHAR(255) NOT NULL,
    bank_name        VARCHAR(255) NOT NULL,
    account_number   VARCHAR(50) NOT NULL,
    ifsc             VARCHAR(20) NOT NULL,
    nickname         VARCHAR(100),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_beneficiaries_user_id ON beneficiaries (user_id);
