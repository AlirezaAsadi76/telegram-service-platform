package rediscatalog

import (
	"context"
	"encoding/json"
	"errors"
	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"

	"github.com/redis/go-redis/v9"
)

func (c *CatalogCache) GetPlatforms(ctx context.Context) ([]smmentity.Platform, bool, error) {
	const op = "rediscatalog.GetPlatforms"

	data, gErr := c.redis.Client().Get(ctx, c.config.PlatformsCacheKey).Bytes()
	if gErr != nil {

		if errors.Is(gErr, redis.Nil) {
			return nil, false, nil
		}
		return nil, false, richerror.New(op, gErr).WithKind(richerror.KindUnexpected).WithMessage(msgerror.CacheReadFailed)
	}

	var platforms []smmentity.Platform
	if err := json.Unmarshal(data, &platforms); err != nil {
		return nil, false, richerror.New(op, err).WithKind(richerror.KindSerializationFailure).WithMessage(msgerror.CacheParseFailed)
	}

	return platforms, true, nil
}
