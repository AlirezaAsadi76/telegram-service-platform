-- +migrate Up

CREATE TABLE ads_plans (
                           id BIGSERIAL PRIMARY KEY,
                           views BIGINT NOT NULL,
                           cpm NUMERIC(18,8) NOT NULL,
                           daily_view_limit BIGINT NOT NULL,
                           active BOOLEAN NOT NULL DEFAULT TRUE,
                           created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                           updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +migrate Down

DROP TABLE ads_plans;