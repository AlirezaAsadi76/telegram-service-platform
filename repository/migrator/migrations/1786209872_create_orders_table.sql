-- +migrate Up

CREATE TABLE orders (
                        id BIGSERIAL PRIMARY KEY,
                        user_id BIGINT NOT NULL,
                        product_type VARCHAR(50) NOT NULL,
                        product_id BIGINT NOT NULL,
                        quantity BIGINT NOT NULL DEFAULT 1,
                        amount NUMERIC(18,8) NOT NULL,
                        currency VARCHAR(20) NOT NULL,
                        status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
                        metadata JSONB,
                        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                        updated_at TIMESTAMP NOT NULL DEFAULT NOW()

);

CREATE INDEX idx_orders_user_id
    ON orders(user_id);

CREATE INDEX idx_orders_status
    ON orders(status);
-- +migrate Down

DROP TABLE orders;