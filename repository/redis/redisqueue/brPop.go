package redisqueue

import (
	"context"
	"errors"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
	"time"

	"github.com/redis/go-redis/v9"
)

func (db DB) BRPop(ctx context.Context, ttl time.Duration, key string) ([]string, error) {
	const Op = "redis_repository.BRPop"

	// BLPOP with 5s timeout — if empty, returns gracefully
	result, bErr := db.adapter.Client().BRPop(ctx, ttl, key).Result()
	if bErr != nil {
		if errors.Is(bErr, redis.Nil) {
			return nil, bErr
		}
		return nil, richerror.New(Op, bErr).WithKind(richerror.KindQueryFailure).WithMessage(msgerror.QueryScanFailed)
	}
	if len(result) < 2 {
		return nil, richerror.New(Op, bErr).WithKind(richerror.KindRedisNil).WithMessage(msgerror.CacheEmpty)
	}

	return result, nil
}
