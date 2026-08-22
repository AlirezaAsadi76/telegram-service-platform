package orderhandler

import (
	"context"
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/params/orderparams"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h *Handler) CanHandle(ctx context.Context, telegramID entity.TelegramId) bool {
	stateResp, err := h.orderFlowService.GetOrderFlow(ctx, orderparams.GetOrderFlowRequest{
		TelegramID: telegramID,
	})
	return err == nil && stateResp != nil && stateResp.Stage != orderentity.OrderFlowStageCompleted
}

func (h *Handler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	h.handleMessage(ctx, b, update)
}
