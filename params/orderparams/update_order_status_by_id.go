package orderparams

import "telegram-service-platform/entity"

type UpdateOrderStatusByIdRequest struct {
	Id     uint64
	Status entity.Status
}

type UpdateOrderStatusByIdResponse struct {
	Id     uint64
	Status entity.Status
}
