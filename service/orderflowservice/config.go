package orderflowservice

import "time"

type Config struct {
	OrderTTL time.Duration `koanf:"order_ttl"`
}
