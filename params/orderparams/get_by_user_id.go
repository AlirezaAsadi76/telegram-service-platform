package orderparams

import "telegram-service-platform/entity/orderentity"

type GetByUserIdRequest struct {
	UserID uint64 `json:"user_id"`
}
type GetByUserIdResponse struct {
	Orders []*orderentity.Order `json:"orders"`
}
