package productservice

import (
	"context"

	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/logger"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"

	"go.uber.org/zap"
)

func (s Service) InvalidateCatalogCache(ctx context.Context, platform smmentity.PlatformType) error {
	const op = "productservice.InvalidateCatalogCache"

	if s.catalogCache == nil {
		return nil
	}

	if err := s.catalogCache.Invalidate(
		ctx,
		platform.String(),
	); err != nil {
		return richerror.New(op, err).
			WithKind(richerror.KindUnexpected).
			WithMessage(msgerror.CacheWriteFailed)
	}

	return nil
}

func (s Service) invalidateCatalogCacheBestEffort(ctx context.Context, platforms ...smmentity.PlatformType) {
	seen := make(
		map[smmentity.PlatformType]struct{},
	)

	for _, platform := range platforms {
		if _, exists := seen[platform]; exists {
			continue
		}

		seen[platform] = struct{}{}

		if err := s.InvalidateCatalogCache(
			ctx,
			platform,
		); err != nil {
			logger.Logger.Warn(
				"catalog cache invalidation failed",
				zap.String("platform", platform.String()),
				zap.Error(err),
			)
		}
	}
}
