-- +migrate Up

ALTER TABLE payments
    ADD COLUMN payment_url TEXT;

-- +migrate Down

ALTER TABLE payments
DROP COLUMN IF EXISTS payment_url;