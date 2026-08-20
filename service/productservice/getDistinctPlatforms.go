package productservice

import (
	"context"
	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/pkg/richerror"
)

func (s Service) GetDistinctPlatforms(ctx context.Context) ([]smmentity.Platform, error) {
	const op = "productservice.GetDistinctPlatforms"

	if platforms, found, err := s.catalogCache.GetPlatforms(ctx); err == nil && found {
		return platforms, nil
	}

	platforms, err := s.repository.SMMMappingGetDistinctPlatforms(ctx)
	if err != nil {
		return nil, richerror.New(op, err).WithKind(richerror.KindQueryFailure)
	}

	if sErr := s.catalogCache.SetPlatforms(ctx, platforms); sErr != nil {
		return nil, richerror.New(op, err).WithKind(richerror.KindQueryFailure)
	}

	return platforms, nil
}
