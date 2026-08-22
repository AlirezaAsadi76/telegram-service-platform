package smmredis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"

	"github.com/redis/go-redis/v9"
)

func (c *SMMCache) GetService(ctx context.Context, id uint64) (*smmentity.SMM, bool, error) {
	const op = "redissmm.GetService"
	key := fmt.Sprintf(serviceKeyPattern, id)

	data, err := c.redis.Client().Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, false, nil
		}
		return nil, false, richerror.New(op, err).
			WithKind(richerror.KindUnexpected).WithMessage(msgerror.CacheReadFailed)
	}

	var service smmentity.SMM
	if err := json.Unmarshal(data, &service); err != nil {
		return nil, false, richerror.New(op, err).
			WithKind(richerror.KindSerializationFailure).WithMessage(msgerror.CacheParseFailed)
	}

	return &service, true, nil
}
