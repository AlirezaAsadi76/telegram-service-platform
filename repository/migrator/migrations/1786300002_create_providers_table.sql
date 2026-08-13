-- +migrate Up
CREATE TABLE IF NOT EXISTS providers (
                                         id BIGSERIAL PRIMARY KEY,
                                         name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('SMM', 'PAYMENT')),
    base_url VARCHAR(500),
    api_key VARCHAR(500),
    config JSONB DEFAULT '{}',
    priority INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

CREATE INDEX IF NOT EXISTS idx_providers_type ON providers(type);
CREATE INDEX IF NOT EXISTS idx_providers_active ON providers(is_active);

-- +migrate Down
DROP TABLE IF EXISTS providers;