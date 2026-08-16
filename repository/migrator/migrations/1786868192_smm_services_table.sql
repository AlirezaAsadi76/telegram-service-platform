-- +migrate Up
CREATE TABLE smm_services (
                              id BIGSERIAL PRIMARY KEY,
                              service_id BIGINT NOT NULL,
                              name TEXT NOT NULL,
                              type VARCHAR(255),
                              rate BIGINT NOT NULL DEFAULT 0,
                              min_quantity BIGINT NOT NULL DEFAULT 100,
                              max_quantity BIGINT NOT NULL DEFAULT 100000,
                              drip_feed BOOLEAN NOT NULL DEFAULT false,
                              refill BOOLEAN NOT NULL DEFAULT false,
                              cancel BOOLEAN NOT NULL DEFAULT false,
                              is_active BOOLEAN NOT NULL DEFAULT true,
                              category VARCHAR(100) NOT NULL,
                              provider_name VARCHAR(100) NOT NULL DEFAULT 'justanotherpanel',
                              created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                              updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_smm_services_service_id ON smm_services(service_id);
CREATE INDEX idx_smm_services_category ON smm_services(category);
CREATE INDEX idx_smm_services_active ON smm_services(is_active);

-- +migrate Down
DROP TABLE IF EXISTS smm_services;