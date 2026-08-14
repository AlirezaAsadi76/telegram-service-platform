package orderfulfillerjob

import (
	"context"
	"time"
)

type RedisRepository interface {
	BRPop(ctx context.Context, ttl time.Duration, key string) ([]string, error)
}
