package orderfulfillmenttesting

import (
	"context"
	"errors"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/entity/providerentity"
	"telegram-service-platform/params/smmparams"
)

type fakeProviderRepository struct {
	provider *providerentity.Provider
}

func (f *fakeProviderRepository) GetActiveByType(
	_ context.Context,
	providerType providerentity.ProviderType,
) ([]*providerentity.Provider, error) {
	if f.provider == nil {
		return nil, nil
	}

	if !f.provider.IsActive {
		return nil, nil
	}

	if f.provider.Type != providerType {
		return nil, nil
	}

	return []*providerentity.Provider{
		f.provider,
	}, nil
}

func (f *fakeProviderRepository) GetByID(
	_ context.Context,
	providerID uint64,
) (*providerentity.Provider, error) {
	if f.provider == nil ||
		f.provider.ID != providerID {
		return nil, errors.New("provider not found")
	}

	return f.provider, nil
}

type fakeSMMProvider struct {
	createResponse smmparams.CreateOrderAdapterResponse
	createErr      error

	status    orderentity.OrderStatus
	statusErr error

	createCalls int
	statusCalls int

	events []string
}

func (f *fakeSMMProvider) Create(
	_ context.Context,
	_ smmparams.CreateOrderAdapterRequest,
) (smmparams.CreateOrderAdapterResponse, error) {
	f.events = append(f.events, "create")
	f.createCalls++

	return f.createResponse, f.createErr
}

func (f *fakeSMMProvider) GetOrderStatus(
	_ context.Context,
	_ string,
) (orderentity.OrderStatus, error) {
	f.events = append(f.events, "get_status")
	f.statusCalls++

	return f.status, f.statusErr
}
