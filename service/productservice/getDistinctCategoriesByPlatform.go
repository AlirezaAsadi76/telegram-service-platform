package productservice

import (
	"context"
	"telegram-service-platform/params/productparams"
	"telegram-service-platform/pkg/richerror"
)

func (s Service) GetDistinctCategoriesByPlatform(ctx context.Context, req productparams.GetDistinctCategoriesByPlatformRequest) (
	productparams.GetDistinctCategoriesByPlatformResponse, error) {
	const op = "productservice.GetDistinctCategoriesByPlatform"

	if categories, found, err := s.catalogCache.GetCategories(ctx, req.Platform.String()); err == nil && found {
		return productparams.GetDistinctCategoriesByPlatformResponse{
			Categories: categories,
		}, nil
	}

	categories, err := s.repository.SMMMappingGetDistinctCategoriesByPlatform(ctx, req.Platform.String())
	if err != nil {
		return productparams.GetDistinctCategoriesByPlatformResponse{}, richerror.New(op, err).WithKind(richerror.KindQueryFailure)
	}

	if sErr := s.catalogCache.SetCategories(ctx, req.Platform.String(), categories); sErr != nil {
		return productparams.GetDistinctCategoriesByPlatformResponse{}, richerror.New(op, err).WithKind(richerror.KindQueryFailure)
	}
	return productparams.GetDistinctCategoriesByPlatformResponse{
		Categories: categories,
	}, nil
}
