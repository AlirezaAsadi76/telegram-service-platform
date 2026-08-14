package redisqueue

import "telegram-service-platform/adapter/redisadapter"

type DB struct {
	adapter *redisadapter.Adapter
}

func New(adapter *redisadapter.Adapter) DB {
	return DB{adapter: adapter}
}
