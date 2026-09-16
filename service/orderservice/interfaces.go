package orderservice

import (
	"context"
	"telegram-service-platform/entity/orderentity"
	"time"
)

type Repository interface {
	Create(ctx context.Context, order *orderentity.Order) error
	GetByID(ctx context.Context, orderID uint64) (*orderentity.Order, error)
	GetByUserID(ctx context.Context, userID uint64) ([]*orderentity.Order, error)
	UpdateStatus(ctx context.Context, id uint64, status orderentity.OrderStatus, externalOrderID string, providerID *uint64) error
	GetByStatus(ctx context.Context, status orderentity.OrderStatus) ([]*orderentity.Order, error)
	GetStalePaid(ctx context.Context, olderThan time.Duration, limit int) ([]*orderentity.Order, error)
	ClaimForProcessing(ctx context.Context, orderID uint64) (bool, error)
	SaveExternalOrder(ctx context.Context, orderID uint64, providerID uint64, externalOrderID string) error
	AssignProvider(ctx context.Context, orderID uint64, providerID uint64) error
}
