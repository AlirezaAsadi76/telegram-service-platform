package rediscatalog

import (
	"context"
	"encoding/json"
	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (c *CatalogCache) SetPlatforms(ctx context.Context, platforms []smmentity.Platform) error {
	const Op = "rediscatalog.setPlatforms"
	data, mErr := json.Marshal(platforms)
	if mErr != nil {
		return richerror.New(Op, mErr).WithKind(richerror.KindSerializationFailure).WithMessage(msgerror.MarshalFailed)
	}
	if err := c.redis.Client().Set(ctx, c.config.PlatformsCacheKey, data, c.config.cacheTTL).Err(); err != nil {
		return richerror.New(Op, mErr).WithKind(richerror.KindMissCatch).WithMessage(msgerror.CacheWriteFailed)
	}
	return nil
}
