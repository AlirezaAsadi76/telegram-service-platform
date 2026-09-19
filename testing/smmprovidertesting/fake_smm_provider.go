package smmprovidertesting

import (
	"context"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/entity/providerentity"
	"telegram-service-platform/params/smmparams"
	"telegram-service-platform/service/smmproviderservice"
)

type fakeSMMProvider struct {
	createResponse smmparams.CreateOrderAdapterResponse
	createErr      error
	createCalls    int
}

func (f *fakeSMMProvider) Create(_ context.Context, _ smmparams.CreateOrderAdapterRequest) (smmparams.CreateOrderAdapterResponse, error) {
	f.createCalls++

	return f.createResponse, f.createErr
}

func (f *fakeSMMProvider) GetOrderStatus(_ context.Context, _ string) (orderentity.OrderStatus, error) {
	return "", nil
}

func newTestService(repo *fakeProviderRepository, adapters map[string]*fakeSMMProvider) *smmproviderservice.Service {
	service := smmproviderservice.New(
		repo,
		smmproviderservice.Config{
			FailureThreshold: 3,
			SuccessThreshold: 1,
		},
	)

	for name, adapter := range adapters {
		service.RegisterProvider(name, adapter)
	}

	return service
}

func newSMMProvider(id uint64, name string) *providerentity.Provider {
	return &providerentity.Provider{
		ID:       id,
		Name:     name,
		Type:     providerentity.ProviderTypeSMM,
		IsActive: true,
	}
}
