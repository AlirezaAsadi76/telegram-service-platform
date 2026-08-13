-- +migrate Up
CREATE TABLE IF NOT EXISTS wallet_transactions (
                                                   id BIGSERIAL PRIMARY KEY,
                                                   wallet_id BIGINT NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('DEPOSIT', 'WITHDRAW', 'REFUND')),
    amount BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'COMPLETED', 'FAILED')),
    reference_id VARCHAR(255),
    idempotency_key VARCHAR(255) UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

CREATE INDEX IF NOT EXISTS idx_wallet_tx_user_id ON wallet_transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_wallet_tx_idempotency ON wallet_transactions(idempotency_key);
CREATE INDEX IF NOT EXISTS idx_wallet_tx_reference ON wallet_transactions(reference_id);

-- +migrate Down
DROP TABLE IF EXISTS wallet_transactions;