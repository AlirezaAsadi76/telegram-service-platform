package redissmm

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

func (c *SMMCache) GetMapping(ctx context.Context, id int64) (*smmentity.SmmMapping, bool, error) {
	const op = "redissmm.GetMapping"
	key := fmt.Sprintf(mappingKeyPattern, id)

	data, err := c.redis.Client().Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, false, nil
		}
		return nil, false, richerror.New(op, err).
			WithKind(richerror.KindUnexpected).WithMessage(msgerror.CacheReadFailed)
	}

	var mapping smmentity.SmmMapping
	if err := json.Unmarshal(data, &mapping); err != nil {
		return nil, false, richerror.New(op, err).
			WithKind(richerror.KindSerializationFailure).WithMessage(msgerror.CacheParseFailed)
	}

	return &mapping, true, nil
}
