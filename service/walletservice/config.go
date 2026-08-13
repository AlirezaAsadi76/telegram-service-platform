package walletservice

import "time"

type Config struct {
	IdempotencyProcessingTTL time.Duration `koanf:"idempotency_processing_ttl"`
	IdempotencyCompletedTTL  time.Duration `koanf:"idempotency_completed_ttl"`
}
