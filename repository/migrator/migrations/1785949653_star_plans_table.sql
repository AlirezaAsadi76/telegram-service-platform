package migrations
-- +migrate Up

CREATE TABLE star_plans (

                            id BIGSERIAL PRIMARY KEY,

                            amount INTEGER NOT NULL,

                            active BOOLEAN NOT NULL DEFAULT TRUE,

                            created_at TIMESTAMP NOT NULL DEFAULT NOW(),

                            updated_at TIMESTAMP NOT NULL DEFAULT NOW()

);


-- +migrate Down

DROP TABLE star_plans;
