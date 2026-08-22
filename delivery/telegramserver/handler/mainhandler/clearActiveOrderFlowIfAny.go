package mainhandler

import (
	"context"
	"telegram-service-platform/entity"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/orderparams"

	"go.uber.org/zap"
)

func (h *Handler) clearActiveOrderFlowIfAny(ctx context.Context, telegramID int64, op string) {
	stateResp, err := h.orderFlowService.GetOrderFlow(ctx, orderparams.GetOrderFlowRequest{TelegramID: entity.TelegramId(telegramID)})
	if err == nil && stateResp != nil {

		_ = h.orderFlowService.AbandonOrderFlow(ctx, orderparams.DeleteOrderFlowRequest{TelegramID: entity.TelegramId(telegramID)}, 0)

		logger.Logger.Info("stale order flow abandoned due to new navigation",
			zap.String("op", op),
			zap.Int64("telegram_id", telegramID),
			zap.String("previous_stage", string(stateResp.Stage)),
		)
	}
}
