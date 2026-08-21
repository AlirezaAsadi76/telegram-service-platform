package orderparams

import (
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/orderentity"
	"time"
)

type SaveOrderFlowRequest struct {
	TelegramID entity.TelegramId
	State      orderentity.OrderFlowState
	TTLMins    time.Duration
}

type SaveOrderFlowResponse struct {
	Saved bool
}

type GetOrderFlowRequest struct {
	TelegramID entity.TelegramId
}

type GetOrderFlowResponse struct {
	State *orderentity.OrderFlowState
	Found bool
}

type DeleteOrderFlowRequest struct {
	TelegramID entity.TelegramId
}

type DeleteOrderFlowResponse struct {
	Deleted bool
}
