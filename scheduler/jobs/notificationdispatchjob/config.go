package notificationdispatchjob

import "time"

type Config struct {
	queueKey     string        `koanf:"queueKey"`
	timeout      time.Duration `koanf:"timeout"`
	PendingLimit int           `koanf:"PendingLimit"`
}
