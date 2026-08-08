package orderparams

import "telegram-service-platform/entity"

type GetOrdersByStatusRequest struct {
	Status entity.Status
}
type GetOrdersByStatusResponse struct {
	OrderInfo []OrderInfo
}
