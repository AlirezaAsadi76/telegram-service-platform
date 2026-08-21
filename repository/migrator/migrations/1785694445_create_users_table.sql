-- +migrate Up

CREATE TABLE users (

                       id BIGSERIAL PRIMARY KEY,
                       telegram_id BIGINT NOT NULL UNIQUE CHECK (telegram_id > 0),
                       username VARCHAR(32) NULL,
                       first_name VARCHAR(255) NULL,
                       last_name VARCHAR(255) NULL,
                       role SMALLINT NOT NULL,
                       last_seen_at TIMESTAMP NOT NULL DEFAULT NOW(),
                       created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMP NOT NULL DEFAULT NOW()

);


-- +migrate Down

DROP TABLE users;