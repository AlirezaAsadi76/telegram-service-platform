package producttesting

import (
	"context"
	"errors"
	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/service/productservice"
)

var _ productservice.Repository = (*fakeProductRepository)(nil)

type fakeProductRepository struct {
	productservice.Repository

	services []smmentity.SMM

	upserted []smmentity.SMM

	failOnServiceID int64

	getAllCalls int
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
