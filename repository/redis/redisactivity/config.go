package redisactivity

import "time"

type Config struct {
	ActivityTTL time.Duration `koanf:"activity_ttl"`
}

const (
	activityKeyPattern  = "user:activity:%d"
	activityScanPattern = "user:activity:*"
	activityKeyPrefix   = "user:activity:"
)
