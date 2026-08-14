package orderfulfillerjob

import (
	"context"
	"time"
)

type RedisRepository interface {
	BRPop(ctx context.Context, key string, ttl time.Duration) ([]string, error)
}
