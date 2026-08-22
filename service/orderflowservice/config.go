package orderflowservice

import "time"

type Config struct {
	OrderTTL time.Duration `json:"order_ttl"`
}
