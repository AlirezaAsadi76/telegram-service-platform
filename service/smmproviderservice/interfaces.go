package smmproviderservice

import (
	"context"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/entity/providerentity"
	"telegram-service-platform/params/orderparams"
)

type SMMProvider interface {
	CreateOrder(ctx context.Context, req orderparams.CreateOrderAdapterRequest) (orderparams.CreateOrderAdapterResponse, error)
	GetOrderStatus(ctx context.Context, externalOrderID string) (orderentity.OrderStatus, error)
}

type Repository interface {
	GetActiveByType(ctx context.Context, providerType providerentity.ProviderType) ([]*providerentity.Provider, error)
}
