package producttesting

import (
	"context"
	"errors"

	"telegram-service-platform/entity"
	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/service/productservice"
)

var _ productservice.Repository = (*fakeProductRepository)(nil)

type fakeProductRepository struct {
	productservice.Repository

	services        []smmentity.SMM
	upserted        []smmentity.SMM
	failOnServiceID int64
	getAllCalls     int

	platforms     []smmentity.Platform
	platformsErr  error
	platformCalls int

	categories    []smmentity.Category
	categoriesErr error
	categoryCalls int

	createMappingErr error
	updateMappingErr error
	getMappingErr    error

	createdMapping  *smmentity.SmmMapping
	updatedMapping  *smmentity.SmmMapping
	existingMapping *smmentity.SmmMapping
}

func (f *fakeProductRepository) SMMServiceCreateOrUpdate(
	_ context.Context,
	service smmentity.SMM,
) error {
	if f.failOnServiceID != 0 &&
		service.Service == f.failOnServiceID {
		return errors.New("database failure")
	}

	f.upserted = append(
		f.upserted,
		service,
	)

	return nil
}

func (f *fakeProductRepository) SMMServiceGetAll(
	_ context.Context,
) ([]smmentity.SMM, error) {
	f.getAllCalls++

	return f.services, nil
}

func (f *fakeProductRepository) SMMMappingGetDistinctPlatforms(
	_ context.Context,
) ([]smmentity.Platform, error) {
	f.platformCalls++

	return f.platforms, f.platformsErr
}

func (f *fakeProductRepository) SMMMappingGetDistinctCategoriesByPlatform(
	_ context.Context,
	_ smmentity.PlatformType,
) ([]smmentity.Category, error) {
	f.categoryCalls++

	return f.categories, f.categoriesErr
}

func (f *fakeProductRepository) SMMMappingCreate(
	_ context.Context,
	mapping *smmentity.SmmMapping,
) error {
	if f.createMappingErr != nil {
		return f.createMappingErr
	}

	mapping.Id = 100

	f.createdMapping = mapping

	return nil
}

func (f *fakeProductRepository) SMMMappingGetByID(
	_ context.Context,
	id int64,
) (*smmentity.SmmMapping, error) {
	if f.getMappingErr != nil {
		return nil, f.getMappingErr
	}

	if f.existingMapping == nil {
		return nil, errors.New("mapping not found")
	}

	if f.existingMapping.Id != id {
		return nil, errors.New("mapping not found")
	}

	return f.existingMapping, nil
}

func (f *fakeProductRepository) SMMMappingUpdate(
	_ context.Context,
	mapping *smmentity.SmmMapping,
) error {
	if f.updateMappingErr != nil {
		return f.updateMappingErr
	}

	f.updatedMapping = mapping

	return nil
}

func newProductServiceWithCatalogCache(
	repository *fakeProductRepository,
	adapter *fakeSMMAdapter,
	cache *fakeCatalogCache,
) *productservice.Service {
	return productservice.New(
		productservice.Config{},
		nil,
		repository,
		cache,
		nil,
		adapter,
	)
}

func newProductService(
	repository *fakeProductRepository,
	adapter *fakeSMMAdapter,
) *productservice.Service {
	return productservice.New(
		productservice.Config{},
		nil,
		repository,
		nil,
		nil,
		adapter,
	)
}

var _ = entity.StarPackage{}
