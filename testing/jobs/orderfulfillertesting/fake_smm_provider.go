package orderfulfillertesting

import (
	"context"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/entity/providerentity"
	"telegram-service-platform/params/smmparams"
)

type fakeSMMProvider struct {
	response smmparams.CreateOrderAdapterResponse
	err      error

	events      []string
	createCall  int
	lastRequest smmparams.CreateOrderAdapterRequest
}

func (f *fakeSMMProvider) Create(_ context.Context, req smmparams.CreateOrderAdapterRequest) (smmparams.CreateOrderAdapterResponse, error) {
	f.createCall++
	f.events = append(f.events, "provider_create")
	f.lastRequest = req

	return f.response, f.err
}

func (f *fakeSMMProvider) GetOrderStatus(_ context.Context, _ string) (orderentity.OrderStatus, error) {
	return "", nil
}

func newProvider(id uint64, name string) *providerentity.Provider {
	return &providerentity.Provider{
		ID:       id,
		Name:     name,
		Type:     providerentity.ProviderTypeSMM,
		IsActive: true,
	}
}
