package orderparams

import "telegram-service-platform/entity/orderentity"

type GetByStatusRequest struct {
	Status orderentity.OrderStatus
}

type GetByStatusResponse struct {
	Orders []*orderentity.Order
}
