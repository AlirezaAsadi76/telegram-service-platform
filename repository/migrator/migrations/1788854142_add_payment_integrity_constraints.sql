-- +migrate Up

CREATE UNIQUE INDEX ux_payments_idempotency_key
    ON payments (idempotency_key)
    WHERE idempotency_key IS NOT NULL;

CREATE UNIQUE INDEX ux_payments_active_order
    ON payments (order_id)
    WHERE status IN (
        'CREATING',
        'PENDING',
        'PROCESSING',
        'UNKNOWN'
    );

-- +migrate Down

DROP INDEX IF EXISTS ux_payments_active_order;
DROP INDEX IF EXISTS ux_payments_idempotency_key;