-- +migrate Up

CREATE TABLE payments (
                          id BIGSERIAL PRIMARY KEY,
                          order_id BIGINT NOT NULL,
                          user_id BIGINT NOT NULL,
                          method VARCHAR(30) NOT NULL,
                          provider VARCHAR(50),
                          amount NUMERIC(18,8) NOT NULL,
                          currency VARCHAR(20) NOT NULL,
                          status VARCHAR(30) NOT NULL,
                          reference VARCHAR(255),
                          metadata JSONB,
                          created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                          updated_at TIMESTAMP NOT NULL DEFAULT NOW()

);

CREATE INDEX idx_payments_order_id
    ON payments(order_id);


CREATE INDEX idx_payments_status
    ON payments(status);

-- +migrate Down

DROP TABLE payments;