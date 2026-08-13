-- +migrate Up
CREATE TABLE IF NOT EXISTS wallets (
                                       id BIGSERIAL PRIMARY KEY,
                                       user_id BIGINT NOT NULL UNIQUE,
                                       balance BIGINT NOT NULL DEFAULT 0,
                                       currency VARCHAR(20) NOT NULL DEFAULT 'IRR',
    version BIGINT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

CREATE INDEX IF NOT EXISTS idx_wallets_user_id ON wallets(user_id);

-- +migrate Down
DROP TABLE IF EXISTS wallets;