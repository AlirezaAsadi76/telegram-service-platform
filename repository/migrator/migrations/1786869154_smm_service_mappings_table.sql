-- +migrate Up
CREATE TABLE smm_service_mapping (
                                     id BIGSERIAL PRIMARY KEY,
                                     smm_service_id BIGINT NOT NULL REFERENCES smm_services(id),
                                     name TEXT NOT NULL,
                                     platform VARCHAR(100) NOT NULL,
                                     category VARCHAR(100) NOT NULL,
                                     description TEXT,
                                     is_active BOOLEAN NOT NULL DEFAULT true,
                                     button_name TEXT,
                                     created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                     updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_smm_mapping_platform ON smm_service_mapping(platform);
CREATE INDEX idx_smm_mapping_category ON smm_service_mapping(category);
CREATE INDEX idx_smm_mapping_active ON smm_service_mapping(is_active);
CREATE INDEX idx_smm_mapping_service ON smm_service_mapping(smm_service_id);

-- +migrate Down
DROP TABLE IF EXISTS smm_service_mapping;