package redisorderflow

import "telegram-service-platform/adapter/redisadapter"

type DB struct {
	redis *redisadapter.Adapter
}

func New(redis *redisadapter.Adapter) *DB {
	return &DB{
		redis: redis,
	}
}
