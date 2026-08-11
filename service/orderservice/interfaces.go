package orderservice

import (
	"context"
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/orderentity"
)

type Repository interface {
	Create(ctx context.Context, order *orderentity.Order) error
	GetByID(ctx context.Context, orderID uint64) (*orderentity.Order, error)
	UpdateStatusByID(ctx context.Context, orderID uint64, status entity.Status) error
	GetByCustomerID(ctx context.Context, customerID uint64) ([]*orderentity.Order, error)
	GetByStatus(ctx context.Context, status entity.Status) ([]*orderentity.Order, error)
}
