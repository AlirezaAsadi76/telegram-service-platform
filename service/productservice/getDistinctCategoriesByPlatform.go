package productservice

import (
	"context"
	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/pkg/richerror"
)

func (s Service) GetDistinctCategoriesByPlatform(ctx context.Context, platform smmentity.PlatformType) ([]smmentity.Category, error) {
	const op = "productservice.GetDistinctCategoriesByPlatform"

	if categories, found, err := s.catalogCache.GetCategories(ctx, platform.String()); err == nil && found {
		return categories, nil
	}

	categories, err := s.repository.SMMMappingGetDistinctCategoriesByPlatform(ctx, platform.String())
	if err != nil {
		return nil, richerror.New(op, err).WithKind(richerror.KindQueryFailure)
	}

	if sErr := s.catalogCache.SetCategories(ctx, platform.String(), categories); sErr != nil {
		return nil, richerror.New(op, err).WithKind(richerror.KindQueryFailure)
	}
	return categories, nil
}
