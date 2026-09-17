-- +migrate Up

CREATE TABLE order_fulfillment_attempts (
                                            id BIGSERIAL PRIMARY KEY,
                                            order_id BIGINT NOT NULL UNIQUE,
                                            provider_id BIGINT NOT NULL,
                                            outcome VARCHAR(20) NOT NULL,
                                            external_order_id VARCHAR(255),
                                            resolved_at TIMESTAMP,
                                            created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                            updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_order_fulfillment_attempts_unresolved
    ON order_fulfillment_attempts(resolved_at)
    WHERE resolved_at IS NULL;

-- +migrate Down

DROP TABLE order_fulfillment_attempts;