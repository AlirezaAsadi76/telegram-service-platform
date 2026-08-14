package notificationdispatchjob

import "time"

type Config struct {
	QueueKey     string        `koanf:"queueKey"`
	Timeout      time.Duration `koanf:"timeout"`
	PendingLimit int           `koanf:"PendingLimit"`
}
