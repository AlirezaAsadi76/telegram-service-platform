package smmredis

import (
	"context"
	"encoding/json"
	"fmt"
	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (c *SMMCache) SetMapping(ctx context.Context, mapping *smmentity.SmmMapping) error {
	const op = "redissmm.SetMapping"
	key := fmt.Sprintf(mappingKeyPattern, mapping.Id)

	data, err := json.Marshal(mapping)
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
