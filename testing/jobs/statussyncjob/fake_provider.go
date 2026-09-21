package statussynctesting

import (
	"context"

	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/entity/providerentity"
	"telegram-service-platform/params/smmparams"
)

type fakeProviderRepository struct {
	provider *providerentity.Provider
}

func (f *fakeProviderRepository) GetActiveByType(
	_ context.Context,
	_ providerentity.ProviderType,
) ([]*providerentity.Provider, error) {
	return nil, nil
}

func (f *fakeProviderRepository) GetByID(
	_ context.Context,
	providerID uint64,
) (*providerentity.Provider, error) {
	if f.provider != nil && f.provider.ID == providerID {
		return f.provider, nil
	}

	return nil, nil
}

type fakeSMMProvider struct {
	status orderentity.OrderStatus
	err    error

	statusCalls int
}

func (f *fakeSMMProvider) Create(
	_ context.Context,
	_ smmparams.CreateOrderAdapterRequest,
) (smmparams.CreateOrderAdapterResponse, error) {
	return smmparams.CreateOrderAdapterResponse{}, nil
}

func (f *fakeSMMProvider) GetOrderStatus(
	_ context.Context,
	_ string,
) (orderentity.OrderStatus, error) {
	f.statusCalls++

	return f.status, f.err
}

func newProvider(
	id uint64,
	name string,
) *providerentity.Provider {
	return &providerentity.Provider{
		ID:       id,
		Name:     name,
		Type:     providerentity.ProviderTypeSMM,
		IsActive: true,
	}
}
