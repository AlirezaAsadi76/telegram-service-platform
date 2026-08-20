package productservice

import (
	"context"
	"telegram-service-platform/entity/smmentity"
)

func (s Service) InvalidateCatalogCache(ctx context.Context, platform smmentity.PlatformType) error {
	return s.catalogCache.Invalidate(ctx, platform.String())
}
