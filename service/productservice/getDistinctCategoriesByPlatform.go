package productservice

import (
	"context"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/productparams"
	"telegram-service-platform/pkg/richerror"

	"go.uber.org/zap"
)

func (s Service) GetDistinctCategoriesByPlatform(ctx context.Context, req productparams.GetDistinctCategoriesByPlatformRequest) (
	productparams.GetDistinctCategoriesByPlatformResponse, error) {
	const op = "productservice.GetDistinctCategoriesByPlatform"

	if s.catalogCache != nil {
		categories, found, cacheErr := s.catalogCache.GetCategories(ctx, req.Platform.String())

		if cacheErr != nil {
			logger.Logger.Warn(
				"catalog categories cache read failed, falling back to database",
				zap.String("op", op),
				zap.String("platform", req.Platform.String()),
				zap.Error(cacheErr),
			)
		} else if found {
			return productparams.GetDistinctCategoriesByPlatformResponse{
				Categories: categories,
			}, nil
		}
	}

	categories, err := s.repository.SMMMappingGetDistinctCategoriesByPlatform(ctx, req.Platform)
	if err != nil {
		return productparams.GetDistinctCategoriesByPlatformResponse{}, richerror.New(op, err).WithKind(richerror.KindQueryFailure)
	}

	if s.catalogCache != nil {
		if cacheErr := s.catalogCache.SetCategories(ctx, req.Platform.String(), categories); cacheErr != nil {
			logger.Logger.Warn(
				"catalog categories cache write failed",
				zap.String("op", op),
				zap.String("platform", req.Platform.String()),
				zap.Error(cacheErr),
			)
		}
	}
	return productparams.GetDistinctCategoriesByPlatformResponse{
		Categories: categories,
	}, nil
}
