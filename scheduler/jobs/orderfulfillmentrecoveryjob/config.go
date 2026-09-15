package orderfulfillmentrecoveryjob

import "time"

type Config struct {
	StaleAfter time.Duration `koanf:"stale_after"`
	BatchSize  int           `koanf:"batch_size"`
}
