package smmproviderservice

import (
	"context"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/entity/providerentity"
	"telegram-service-platform/params/smmparams"
)

type SMMProvider interface {
	Create(ctx context.Context, req smmparams.CreateOrderAdapterRequest) (smmparams.CreateOrderAdapterResponse, error)
	GetOrderStatus(ctx context.Context, externalOrderID string) (orderentity.OrderStatus, error)
}

type ProviderRepository interface {
	GetActiveByType(ctx context.Context, providerType providerentity.ProviderType) ([]*providerentity.Provider, error)
	GetByID(ctx context.Context, providerID uint64) (*providerentity.Provider, error)
}
