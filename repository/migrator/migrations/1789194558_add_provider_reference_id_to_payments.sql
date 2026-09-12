-- +migrate Up

ALTER TABLE payments
    ADD COLUMN provider_reference_id VARCHAR(255);

CREATE INDEX idx_payments_provider_reference_id
    ON payments(provider_reference_id);

-- +migrate Down

DROP INDEX IF EXISTS idx_payments_provider_reference_id;

ALTER TABLE payments
DROP COLUMN provider_reference_id;