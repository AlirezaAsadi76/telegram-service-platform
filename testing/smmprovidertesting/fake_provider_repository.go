package smmprovidertesting

import (
	"context"
	"errors"
	"telegram-service-platform/entity/providerentity"
)

type fakeProviderRepository struct {
	providers []*providerentity.Provider
}

func (f *fakeProviderRepository) GetActiveByType(_ context.Context, providerType providerentity.ProviderType) ([]*providerentity.Provider, error) {
	result := make([]*providerentity.Provider, 0)

	for _, provider := range f.providers {
		if provider.Type == providerType && provider.IsActive {
			result = append(result, provider)
		}
	}

	return result, nil
}

func (f *fakeProviderRepository) GetByID(_ context.Context, providerID uint64) (*providerentity.Provider, error) {
	for _, provider := range f.providers {
		if provider.ID == providerID {
			return provider, nil
		}
	}

	return nil, errors.New("provider not found")
}
