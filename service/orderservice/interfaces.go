package orderservice

import (
	"context"
	"telegram-service-platform/entity/orderentity"
)

type Repository interface {
	Create(ctx context.Context, order *orderentity.Order) error
	GetByID(ctx context.Context, orderID uint64) (*orderentity.Order, error)
	UpdateStatus(ctx context.Context, id uint64, status orderentity.OrderStatus, externalOrderID string, providerID *uint64) error
	GetByStatus(ctx context.Context, status orderentity.OrderStatus) ([]*orderentity.Order, error)
}
