package rediscatalog

import (
	"context"
	"encoding/json"
	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (c *CatalogCache) GetPlatforms(ctx context.Context) ([]smmentity.Platform, bool, error) {
	const Op = "rediscatalog.getPlatforms"
	data, gErr := c.redis.Client().Get(ctx, c.config.PlatformsCacheKey).Bytes()
	if gErr != nil {
		return nil, false, richerror.New(Op, gErr).WithKind(richerror.KindMissCatch).WithMessage(msgerror.CacheReadFailed)
	}

	var platforms []smmentity.Platform
	if err := json.Unmarshal(data, &platforms); err != nil {
		return nil, false, richerror.New(Op, gErr).WithKind(richerror.KindSerializationFailure).WithMessage(msgerror.CacheParseFailed)
	}
	return platforms, true, nil
}
