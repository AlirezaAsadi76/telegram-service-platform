package orderparams

import "telegram-service-platform/entity/orderentity"

type CreateOrderAdapterRequest struct {
	ServiceID string
	Link      string
	Quantity  int64
}

type CreateOrderAdapterResponse struct {
	ExternalOrderID string
	Status          orderentity.OrderStatus
}
