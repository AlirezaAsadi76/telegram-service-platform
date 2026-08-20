package rediscatalog

import (
	"context"
	"fmt"
)

func (c *CatalogCache) Invalidate(ctx context.Context, platform string) error {
	_ = c.redis.Client().Del(ctx, fmt.Sprintf(c.config.categoriesCacheKey, platform)).Err()
	_ = c.redis.Client().Del(ctx, c.config.PlatformsCacheKey).Err()
	return nil
}
