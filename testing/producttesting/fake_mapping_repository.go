package producttesting

import (
	"context"
	"errors"
	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/service/productservice"
)

type mappingResolutionRepository struct {
	productservice.Repository

	mapping *smmentity.SmmMapping
	service *smmentity.SMM

	mappingErr error
	serviceErr error
}

func (f *mappingResolutionRepository) SMMMappingGetByID(
	_ context.Context,
	id int64,
) (*smmentity.SmmMapping, error) {
	if f.mappingErr != nil {
		return nil, f.mappingErr
	}

	if f.mapping == nil ||
		f.mapping.Id != id {
		return nil, errors.New("mapping not found")
	}

	return f.mapping, nil
}

func (f *mappingResolutionRepository) SMMServiceGetByD(
	_ context.Context,
	id int64,
) (*smmentity.SMM, error) {
	if f.serviceErr != nil {
		return nil, f.serviceErr
	}

	if f.service == nil ||
		f.service.Id != id {
		return nil, errors.New("smm service not found")
	}

	return f.service, nil
}
