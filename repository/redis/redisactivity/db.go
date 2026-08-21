package redisactivity

import (
	"telegram-service-platform/adapter/redisadapter"
	"time"
)

type ActivityTracker struct {
	redis *redisadapter.Adapter
	ttl   time.Duration
}

func New(redis *redisadapter.Adapter, cfg Config) *ActivityTracker {
	return &ActivityTracker{
		redis: redis,
		ttl:   cfg.ActivityTTL,
	}
}
