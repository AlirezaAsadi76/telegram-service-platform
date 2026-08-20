package rediscatalog

import (
	"context"
	"encoding/json"
	"fmt"
	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (c *CatalogCache) GetCategories(ctx context.Context, platform string) ([]smmentity.Category, bool, error) {
	const Op = "rediscatalog.getCategories"
	key := fmt.Sprintf(c.config.categoriesCacheKey, platform)
	data, gErr := c.redis.Client().Get(ctx, key).Bytes()
	if gErr != nil {
		return nil, false, richerror.New(Op, gErr).WithKind(richerror.KindMissCatch).WithMessage(msgerror.CacheReadFailed)
	}

	var categories []smmentity.Category
	if err := json.Unmarshal(data, &categories); err != nil {
		return nil, false, richerror.New(Op, gErr).WithKind(richerror.KindSerializationFailure).WithMessage(msgerror.CacheParseFailed)
	}
	return categories, true, nil
}
