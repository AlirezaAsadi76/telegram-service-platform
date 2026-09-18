package orderfulfillmentrecoveryjob

import "time"

type Config struct {
	StaleAfter          time.Duration `koanf:"stale_after"`
	ReconciliationAfter time.Duration `koanf:"reconciliation_after"`
	BatchSize           int           `koanf:"batch_size"`
}
