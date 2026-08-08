package orderservice

import (
	"context"
	"telegram-service-platform/entity"
)

type Repository interface {
	Create(ctx context.Context, order *entity.Order) error
	GetByID(ctx context.Context, orderID uint64) (*entity.Order, error)
	UpdateStatusByID(ctx context.Context, orderID uint64, status entity.Status) error
	GetByCustomerID(ctx context.Context, customerID uint64) ([]*entity.Order, error)
	GetByStatus(ctx context.Context, status entity.Status) ([]*entity.Order, error)
}
