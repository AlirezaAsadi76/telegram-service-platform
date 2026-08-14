package orderfulfillerjob

import "time"

type Config struct {
	QueueKey string        `koanf:"queue_key"`
	Timeout  time.Duration `koanf:"timeout"`
}
