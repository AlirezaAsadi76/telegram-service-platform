package rediscatalog

import (
	"context"
	"fmt"

	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (c *CatalogCache) Invalidate(ctx context.Context, platform string) error {
	const op = "rediscatalog.Invalidate"

	categoryKey := fmt.Sprintf(c.config.categoriesCacheKey, platform)

	var firstErr error

	if err := c.redis.Client().Del(ctx, categoryKey).Err(); err != nil {
		firstErr = err
	}

	if err := c.redis.Client().Del(ctx, c.config.PlatformsCacheKey).Err(); err != nil && firstErr == nil {
		firstErr = err
	}

	if firstErr != nil {
		return richerror.New(op, firstErr).
			WithKind(richerror.KindUnexpected).
			WithMessage(msgerror.CacheWriteFailed)
	}

	return nil
}
