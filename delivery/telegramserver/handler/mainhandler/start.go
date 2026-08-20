package mainhandler

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"

	"telegram-service-platform/delivery/telegramserver/keyboard"
	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/productparams"
	"telegram-service-platform/pkg/metrics"
	"telegram-service-platform/pkg/richerror"
)

func (h *Handler) start(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "mainhandler.start"

	if update.Message == nil {
		return
	}
	chatID := update.Message.Chat.ID

	metrics.OrdersTotal.WithLabelValues("start_command", "triggered").Inc()

	platformsResp, err := h.productService.GetDistinctPlatforms(ctx)
	if err != nil {
		h.handleError(ctx, b, chatID, op, err)
		return
	}

	categoriesResp, gErr := h.productService.GetDistinctCategoriesByPlatform(ctx, productparams.GetDistinctCategoriesByPlatformRequest{
		Platform: smmentity.TelegramPlatform,
	})
	if gErr != nil {
		h.handleError(ctx, b, chatID, op, gErr)
		return
	}

	welcomeMsg := fmt.Sprintf(
		"سلام %s! 👋\nبه پلتفرم خدمات SMM خوش آمدید.\nلطفاً از منوی زیر انتخاب کنید:",
		update.Message.From.FirstName,
	)

	inlineKeyboard := keyboard.MainMenu(platformsResp.Platforms, categoriesResp.Categories)
	if sendErr := h.messenger.Send(ctx, b, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        welcomeMsg,
		ReplyMarkup: inlineKeyboard,
	}); sendErr != nil {
		logger.Logger.Error("failed to send welcome message", zap.String("op", op), zap.Error(sendErr))
	}

	replyKeyboard := keyboard.ReplyMainMenu()
	if sendErr := h.messenger.Send(ctx, b, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        "📌 دکمه‌های دسترسی سریع پایین صفحه همیشه در دسترس شما هستند:",
		ReplyMarkup: replyKeyboard,
	}); sendErr != nil {
		logger.Logger.Error("failed to send reply keyboard", zap.String("op", op), zap.Error(sendErr))
	}
}

func (h *Handler) showMainMenu(ctx context.Context, b *bot.Bot, update *models.Update) {
	h.start(ctx, b, update)
}

func (h *Handler) handleError(ctx context.Context, b *bot.Bot, chatID any, op string, err error) {

	if richErr, ok := errors.AsType[*richerror.RichError](err); ok {
		logger.Logger.Error("business logic error in handler",
			zap.String("op", op),
			zap.String("kind", string(richErr.Kind())),
			zap.String("message", richErr.Message()),
			zap.Error(richErr.Unwrap()),
		)
	} else {
		logger.Logger.Error("unexpected system error in handler",
			zap.String("op", op),
			zap.Error(err),
		)
	}

	_ = h.messenger.Send(ctx, b, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "⚠️ خطایی در پردازش درخواست شما رخ داده است. لطفاً چند لحظه بعد دوباره تلاش کنید.",
	})
}
