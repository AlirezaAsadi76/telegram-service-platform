package scheduler

import "time"

type Config struct {
	CurrencyRefreshInterval time.Duration `koanf:"currency_refresh_interval"`
	StarsRefreshInterval    time.Duration `koanf:"stars_refresh_interval"`
	PremiumRefreshInterval  time.Duration `koanf:"premium_refresh_interval"`

	PaymentVerifyInterval time.Duration `koanf:"payment_verify_interval"`
	StatusSyncInterval    time.Duration `koanf:"status_sync_interval"`
	PaymentExpiryInterval time.Duration `koanf:"payment_expiry_interval"`

	QueueConsumerInterval            time.Duration `koanf:"queue_consumer_interval"`
	OrderFulfillmentRecoveryInterval time.Duration `koanf:"order_fulfillment_recovery_interval"`

	SmmValidationInterval time.Duration `koanf:"smm_validation_interval"`

	userActivitySyncInterval time.Duration `koanf:"user_activity_sync_interval"`
}
