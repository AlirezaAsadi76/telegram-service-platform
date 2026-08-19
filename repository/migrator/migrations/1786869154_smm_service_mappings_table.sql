-- +migrate Up
CREATE TABLE smm_service_mappings (
                                     id BIGSERIAL PRIMARY KEY,
                                     smm_service_id BIGINT NOT NULL REFERENCES smm_services(id),
                                     name TEXT NOT NULL,
                                     platform VARCHAR(100) NOT NULL,
                                     category VARCHAR(100) NOT NULL,
                                     description TEXT,
                                     is_active BOOLEAN NOT NULL DEFAULT true,
                                     button_name TEXT,
                                     sort_order INT NOT NULL DEFAULT 0
                                     created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                     updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_smm_mappings_platform ON smm_service_mappings(platform);
CREATE INDEX idx_smm_mappings_category ON smm_service_mappings(category);
CREATE INDEX idx_smm_mappings_active ON smm_service_mappings(is_active);
CREATE INDEX idx_smm_mappings_service ON smm_service_mappings(smm_service_id);

-- +migrate Down
DROP TABLE IF EXISTS smm_service_mappings;