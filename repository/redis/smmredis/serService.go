package smmredis

import (
	"context"
	"encoding/json"
	"fmt"
	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (c *SMMCache) SetService(ctx context.Context, service *smmentity.SMM) error {
	const op = "redissmm.SetService"
	key := fmt.Sprintf(serviceKeyPattern, service.Service)

	data, err := json.Marshal(service)
	if err != nil {
		return richerror.New(op, err).
			WithKind(richerror.KindSerializationFailure).WithMessage(msgerror.MarshalFailed)
	}

	if err := c.redis.Client().Set(ctx, key, data, c.cfg.SMMTTL).Err(); err != nil {
		return richerror.New(op, err).
			WithKind(richerror.KindUnexpected).WithMessage(msgerror.CacheWriteFailed)
	}
	return nil
}
