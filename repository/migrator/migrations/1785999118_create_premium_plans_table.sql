-- +migrate Up

CREATE TABLE premium_plans (
                               id BIGSERIAL PRIMARY KEY,
                               duration SMALLINT NOT NULL,
                               active BOOLEAN NOT NULL DEFAULT TRUE,
                               created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                               updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +migrate Down

DROP TABLE premium_plans;