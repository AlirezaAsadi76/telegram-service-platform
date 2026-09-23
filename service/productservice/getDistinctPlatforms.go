package productservice

import (
	"context"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/productparams"
	"telegram-service-platform/pkg/richerror"

	"go.uber.org/zap"
)

func (s Service) GetDistinctPlatforms(ctx context.Context) (productparams.GetDistinctPlatformsResponse, error) {
	const op = "productservice.GetDistinctPlatforms"

	if s.catalogCache != nil {
		platforms, found, cacheErr := s.catalogCache.GetPlatforms(ctx)

		if cacheErr != nil {
			logger.Logger.Warn(
				"catalog platforms cache read failed, falling back to database",
				zap.String("op", op),
				zap.Error(cacheErr),
			)
		} else if found {
			return productparams.GetDistinctPlatformsResponse{
				Platforms: platforms,
			}, nil
		}
	}

	platforms, err := s.repository.SMMMappingGetDistinctPlatforms(ctx)
	if err != nil {
		return productparams.GetDistinctPlatformsResponse{}, richerror.New(op, err).WithKind(richerror.KindQueryFailure)
	}

	if cacheErr := s.catalogCache.SetPlatforms(ctx, platforms); cacheErr != nil {
		logger.Logger.Warn(
			"catalog platforms cache write failed",
			zap.String("op", op),
			zap.Error(cacheErr),
		)
	}
	return productparams.GetDistinctPlatformsResponse{
		Platforms: platforms,
	}, nil
}
