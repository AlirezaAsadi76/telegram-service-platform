package orderparams

import "telegram-service-platform/entity/orderentity"

type UpdateStatusRequest struct {
	OrderID         uint64
	Status          orderentity.OrderStatus
	ExternalOrderID string
	ProviderID      *uint64
}
