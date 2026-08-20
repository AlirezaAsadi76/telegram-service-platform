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

func (h *Handler) selectCategory(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "mainhandler.selectCategory"

	if update.CallbackQuery == nil || update.CallbackQuery.Message.Message == nil {
		return
	}

	chatID := update.CallbackQuery.Message.Message.Chat.ID

	data, ccErr := h.callbackQueryData(update.CallbackQuery.Data, CategorySplitMode)
	platformName := data[platformSplitMode]
	categoryName := data[CategorySplitMode]
	if ccErr != nil {
		h.handleError(ctx, b, chatID, op, fmt.Errorf("داده نامعتبر"))
		return
	}

	// دریافت سرویس‌های این دسته‌بندی از سرویس
	servicesResp, err := h.productService.GetSMMMappingsByPlatformCategory(ctx, productparams.GetSmmMappingByPlatformCategoryRequest{
		Platform: smmentity.PlatformType(platformName),
		Category: smmentity.CategoryType(categoryName),
	})
	if err != nil {
		logger.Logger.Error("failed to get services",
			zap.String("op", op),
			zap.String("platform", platformName),
			zap.String("category", categoryName),
			zap.Error(err),
		)
		h.handleError(ctx, b, chatID, op, err)
		return
	}

	if len(servicesResp.SmmMapping) == 0 {
		h.handleError(ctx, b, chatID, op, fmt.Errorf("هیچ سرویسی برای این دسته‌بندی یافت نشد"))
		return
	}

	message := fmt.Sprintf(
		"%s %s\nلطفاً سرویس مورد نظر خود را انتخاب کنید:",
		helpers.GetCategoryIcon(categoryName),
		helpers.GetCategoryDisplayName(categoryName),
	)

	keyboardMarkup := keyboard.ServiceMenu(platformName, categoryName, servicesResp.SmmMapping)

	if editErr := h.messenger.Edit(ctx, b, &bot.EditMessageTextParams{
		ChatID:      chatID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		Text:        message,
		ReplyMarkup: keyboardMarkup,
	}); editErr != nil {
		logger.Logger.Error("failed to edit message", zap.Error(editErr), zap.String("op", op))
	}
}
