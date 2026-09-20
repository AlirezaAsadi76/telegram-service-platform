package orderfulfillertesting

import (
	"context"
	"telegram-service-platform/entity/providerentity"
)

type fakeProviderRepository struct {
	provider *providerentity.Provider
}

func (f *fakeProviderRepository) GetActiveByType(_ context.Context, providerType providerentity.ProviderType) ([]*providerentity.Provider, error) {
	if f.provider == nil ||
		!f.provider.IsActive ||
		f.provider.Type != providerType {
		return nil, nil
	}

	return []*providerentity.Provider{
		f.provider,
	}, nil
}

func (f *fakeProviderRepository) GetByID(_ context.Context, providerID uint64) (*providerentity.Provider, error) {
	if f.provider != nil &&
		f.provider.ID == providerID {
		return f.provider, nil
	}

	return nil, nil
}
