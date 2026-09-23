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
	const op = "rediscatalog.SetCategories"

	key := fmt.Sprintf(
		c.config.categoriesCacheKey,
		platform,
	)

	data, err := json.Marshal(categories)
	if err != nil {
		return richerror.New(op, err).
			WithKind(richerror.KindSerializationFailure).
			WithMessage(msgerror.MarshalFailed)
	}

	if err := c.redis.Client().Set(ctx, key, data, c.config.cacheTTL).Err(); err != nil {
		return richerror.New(op, err).
			WithKind(richerror.KindUnexpected).
			WithMessage(msgerror.CacheWriteFailed)
	}

	return nil
}
