package rediscatalog

import (
	"context"
	"encoding/json"
	"fmt"
	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (c *CatalogCache) SetCategories(ctx context.Context, platform string, categories []smmentity.Category) error {
	const Op = "rediscatalog.SetCategories"
	key := fmt.Sprintf(c.config.categoriesCacheKey, platform)
	data, mErr := json.Marshal(categories)
	if mErr != nil {
		return richerror.New(Op, mErr).WithKind(richerror.KindSerializationFailure).WithMessage(msgerror.MarshalFailed)
	}

	if err := c.redis.Client().Set(ctx, key, data, c.config.cacheTTL).Err(); err != nil {
		return richerror.New(Op, mErr).WithKind(richerror.KindMissCatch).WithMessage(msgerror.CacheWriteFailed)
	}

	return nil
}
