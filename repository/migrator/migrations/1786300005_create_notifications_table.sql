-- +migrate Up
CREATE TABLE IF NOT EXISTS notifications (
                                            id BIGSERIAL PRIMARY KEY,
                                            user_id BIGINT NOT NULL,
                                            type VARCHAR(50) NOT NULL,
                                            status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
                                            payload JSONB,
                                            retry_count INT NOT NULL DEFAULT 0,
                                            created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                            sent_at TIMESTAMP
    );

CREATE INDEX idx_notifications_status ON notifications(status);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);

-- +migrate Down
DROP TABLE IF EXISTS notifications;