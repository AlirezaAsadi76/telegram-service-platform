package redissmm

import "telegram-service-platform/adapter/redisadapter"

type SMMCache struct {
	redis *redisadapter.Adapter
	cfg   Config
}

func New(redis *redisadapter.Adapter, cfg Config) *SMMCache {
	return &SMMCache{
		redis: redis,
		cfg:   cfg,
	}
}
