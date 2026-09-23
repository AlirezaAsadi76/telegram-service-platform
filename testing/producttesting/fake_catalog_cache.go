package producttesting

import (
	"context"

	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/service/productservice"
)

var _ productservice.CatalogCache = (*fakeCatalogCache)(nil)

type fakeCatalogCache struct {
	platforms      []smmentity.Platform
	platformsFound bool

	categories      []smmentity.Category
	categoriesFound bool

	getPlatformsErr  error
	setPlatformsErr  error
	getCategoriesErr error
	setCategoriesErr error

	invalidatedPlatforms []string
	invalidateErr        error

	setPlatformsCalls  int
	setCategoriesCalls int
	getPlatformsCalls  int
	getCategoriesCalls int
}

func (f *fakeCatalogCache) GetPlatforms(
	_ context.Context,
) ([]smmentity.Platform, bool, error) {
	f.getPlatformsCalls++

	return f.platforms, f.platformsFound, f.getPlatformsErr
}

func (f *fakeCatalogCache) SetPlatforms(
	_ context.Context,
	platforms []smmentity.Platform,
) error {
	f.setPlatformsCalls++

	f.platforms = platforms

	return f.setPlatformsErr
}

func (f *fakeCatalogCache) GetCategories(
	_ context.Context,
	_ string,
) ([]smmentity.Category, bool, error) {
	f.getCategoriesCalls++

	return f.categories, f.categoriesFound, f.getCategoriesErr
}

func (f *fakeCatalogCache) SetCategories(
	_ context.Context,
	_ string,
	categories []smmentity.Category,
) error {
	f.setCategoriesCalls++

	f.categories = categories

	return f.setCategoriesErr
}

func (f *fakeCatalogCache) Invalidate(
	_ context.Context,
	platform string,
) error {
	f.invalidatedPlatforms = append(
		f.invalidatedPlatforms,
		platform,
	)

	return f.invalidateErr
}
