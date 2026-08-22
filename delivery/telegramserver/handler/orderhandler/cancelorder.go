package orderhandler

import (
	"context"
	"telegram-service-platform/entity"
	"telegram-service-platform/params/orderparams"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"

	"telegram-service-platform/logger"
)

func (h *Handler) cancelOrder(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "orderhandler.cancelOrder"

	if update.CallbackQuery == nil || update.CallbackQuery.Message.Message == nil {
		return
	}

	chatID := update.CallbackQuery.Message.Message.ID
	telegramID := update.CallbackQuery.From.ID

	err := h.orderFlowService.AbandonOrderFlow(ctx, orderparams.DeleteOrderFlowRequest{
		TelegramID: entity.TelegramId(telegramID),
	}, 0)

	if err != nil {
		logger.Logger.Error("failed to abandon order flow",
			zap.String("op", op),
			zap.Int64("telegram_id", telegramID),
			zap.Error(err),
		)
	}

	_ = h.messenger.Send(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "❌ سفارش شما با موفقیت لغو شد.\nمی‌توانید از منوی اصلی سرویس جدیدی انتخاب کنید.",
	})
}
