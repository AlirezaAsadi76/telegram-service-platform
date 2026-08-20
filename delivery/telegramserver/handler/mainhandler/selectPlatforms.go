package mainhandler

import (
	"context"
	"fmt"
	"telegram-service-platform/pkg/helpers"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"

	"telegram-service-platform/delivery/telegramserver/keyboard"
	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/productparams"
)

func (h *Handler) selectPlatform(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "mainhandler.selectPlatform"

	if update.CallbackQuery == nil || update.CallbackQuery.Message.Message == nil {
		return
	}
	chatID := update.CallbackQuery.Message.Message.Chat.ID

	data, cErr := h.callbackQueryData(update.CallbackQuery.Data, platformSplitMode)
	platformName := data[platformSplitMode]
	if cErr != nil {
		h.handleError(ctx, chatID, op, fmt.Errorf("داده نامعتبر"))
		return
	}

	categoriesResp, err := h.productService.GetDistinctCategoriesByPlatform(ctx, productparams.GetDistinctCategoriesByPlatformRequest{
		Platform: smmentity.PlatformType(platformName),
	})
	if err != nil {
		h.handleError(ctx, chatID, op, err)
		return
	}

	if len(categoriesResp.Categories) == 0 {

		_ = h.messenger.Send(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "⚠️ هیچ دسته‌بندی فعالی برای این پلتفرم یافت نشد.",
		})
		return
	}

	message := fmt.Sprintf("📱 %s\nلطفاً دسته‌بندی مورد نظر خود را انتخاب کنید:", helpers.GetPlatformDisplayName(platformName))
	keyboardMarkup := keyboard.CategoryMenu(platformName, categoriesResp.Categories)

	if editErr := h.messenger.Edit(ctx, &bot.EditMessageTextParams{
		ChatID:      chatID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		Text:        message,
		ReplyMarkup: keyboardMarkup,
	}); editErr != nil {
		logger.Logger.Error("failed to edit message", zap.String("op", op), zap.Error(editErr))
	}
}
