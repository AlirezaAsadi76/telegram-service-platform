package rediscatalog

import (
	"telegram-service-platform/adapter/redisadapter"
)

type CatalogCache struct {
	redis  *redisadapter.Adapter
	config Config
}

func New(redis *redisadapter.Adapter, cfg Config) *CatalogCache {
	return &CatalogCache{
		redis:  redis,
		config: cfg,
	}
}
