-- +migrate Up

CREATE TABLE payments (
                          id BIGSERIAL PRIMARY KEY,
                          order_id BIGINT NOT NULL,
                          user_id BIGINT NOT NULL,
                          method VARCHAR(30) NOT NULL,
                          amount NUMERIC(18,8) NOT NULL,
                          currency VARCHAR(20) NOT NULL,
                          status VARCHAR(30) NOT NULL,
                          external_id VARCHAR(255),
                          idempotency_key varchar(255),
                          callback_data JSONB,
                          expired_at TIMESTAMP,
                          created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                          updated_at TIMESTAMP NOT NULL DEFAULT NOW()

);

CREATE INDEX idx_payments_order_id
    ON payments(order_id);


CREATE INDEX idx_payments_external_id
    ON payments(external_id);


CREATE INDEX idx_payments_status
    ON payments(status);

-- +migrate Down

DROP TABLE payments;