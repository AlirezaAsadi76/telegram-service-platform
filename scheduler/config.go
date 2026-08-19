package scheduler

import "time"

type Config struct {
	CurrencyRefreshInterval time.Duration `koanf:"currency_refresh_interval"`
	StarsRefreshInterval    time.Duration `koanf:"stars_refresh_interval"`
	PremiumRefreshInterval  time.Duration `koanf:"premium_refresh_interval"`

	// Phase 2 — Cron Jobs
	PaymentVerifyInterval time.Duration `koanf:"payment_verify_interval"`
	StatusSyncInterval    time.Duration `koanf:"status_sync_interval"`
	PaymentExpiryInterval time.Duration `koanf:"payment_expiry_interval"`

	// Phase 2 — Queue Consumers (short interval for continuous polling)
	QueueConsumerInterval time.Duration `koanf:"queue_consumer_interval"`
	// validation services
	SmmValidationInterval time.Duration `koanf:"smm_validation_interval"`
}
