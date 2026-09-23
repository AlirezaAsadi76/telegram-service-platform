package rediscatalog

import (
	"context"
	"encoding/json"

	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (c *CatalogCache) SetPlatforms(ctx context.Context, platforms []smmentity.Platform) error {
	const op = "rediscatalog.SetPlatforms"

	data, err := json.Marshal(platforms)
	if err != nil {
		return richerror.New(op, err).
			WithKind(richerror.KindSerializationFailure).
			WithMessage(msgerror.MarshalFailed)
	}

	if err := c.redis.Client().Set(ctx, c.config.PlatformsCacheKey, data, c.config.cacheTTL).Err(); err != nil {
		return richerror.New(op, err).
			WithKind(richerror.KindUnexpected).
			WithMessage(msgerror.CacheWriteFailed)
	}

	return nil
}
