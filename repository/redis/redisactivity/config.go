package redisactivity

import "time"

type Config struct {
	TTL time.Duration `koanf:"ttl"`
}

const (
	activityKeyPattern  = "user:activity:%d"
	activityScanPattern = "user:activity:*"
	activityKeyPrefix   = "user:activity:"
)
