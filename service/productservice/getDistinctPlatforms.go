package productservice

import (
	"context"
	"telegram-service-platform/params/productparams"
	"telegram-service-platform/pkg/richerror"
)

func (s Service) GetDistinctPlatforms(ctx context.Context) (productparams.GetDistinctPlatformsResponse, error) {
	const op = "productservice.GetDistinctPlatforms"

	if platforms, found, err := s.catalogCache.GetPlatforms(ctx); err == nil && found {
		return productparams.GetDistinctPlatformsResponse{
			Platforms: platforms,
		}, nil
	}

	platforms, err := s.repository.SMMMappingGetDistinctPlatforms(ctx)
	if err != nil {
		return productparams.GetDistinctPlatformsResponse{}, richerror.New(op, err).WithKind(richerror.KindQueryFailure)
	}

	if sErr := s.catalogCache.SetPlatforms(ctx, platforms); sErr != nil {
		return productparams.GetDistinctPlatformsResponse{}, richerror.New(op, err).WithKind(richerror.KindQueryFailure)
	}

	return productparams.GetDistinctPlatformsResponse{
		Platforms: platforms,
	}, nil
}
