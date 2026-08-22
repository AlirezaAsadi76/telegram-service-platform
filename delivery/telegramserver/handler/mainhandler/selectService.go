package mainhandler

import (
	"context"
	"fmt"
	"strconv"
	"telegram-service-platform/pkg/helpers"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"

	"telegram-service-platform/logger"
	"telegram-service-platform/params/productparams"
)

func (h *Handler) selectService(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "mainhandler.selectService"

	if update.CallbackQuery == nil || update.CallbackQuery.Message.Message == nil {
		return
	}

	chatID := update.CallbackQuery.Message.Message.Chat.ID
	messageID := update.CallbackQuery.Message.Message.ID
	h.clearActiveOrderFlowIfAny(ctx, update.CallbackQuery.From.ID, op)

	data, ccErr := h.callbackQueryData(update.CallbackQuery.Data, CategorySplitMode)
	platformName := data[platformSplitMode]
	categoryName := data[CategorySplitMode]
	serviceIDStr := data[ServiceSplitMode]
	if ccErr != nil {
		h.handleError(ctx, chatID, op, fmt.Errorf("داده نامعتبر"))
		return
	}

	serviceID, err := strconv.ParseInt(serviceIDStr, 10, 64)
	if err != nil {
		logger.Logger.Error("invalid service ID",
			zap.String("op", op),
			zap.String("serviceID", serviceIDStr),
		)
		h.handleError(ctx, chatID, op, err)
		return
	}

	mappingResp, gErr := h.productService.GetSMMMappingByID(ctx, productparams.GetSmmMappingByIDRequest{
		Id: serviceID,
	})
	if gErr != nil {
		logger.Logger.Error("failed to get service details",
			zap.String("op", op),
			zap.Int64("serviceID", serviceID),
			zap.Error(gErr),
		)
		h.handleError(ctx, chatID, op, err)
		return
	}

	message := fmt.Sprintf(
		"✅ سرویس انتخاب شد:\n\n"+
			"📱 پلتفرم: %s\n"+
			"%s دسته‌بندی: %s\n"+
			"📦 سرویس: %s\n\n"+
			"لطفاً لینک مورد نظر خود را ارسال کنید:",
		helpers.GetPlatformDisplayName(platformName),
		helpers.GetCategoryIcon(categoryName),
		helpers.GetCategoryDisplayName(categoryName),
		mappingResp.SmmMapping.ButtonName,
	)

	if sendErr := h.messenger.Edit(ctx, &bot.EditMessageTextParams{
		ChatID:    chatID,
		Text:      message,
		MessageID: messageID,
	}); sendErr != nil {
		logger.Logger.Error("failed to send service details", zap.Error(sendErr), zap.String("op", op))
	}

	// TODO: در فاز بعدی، state کاربر را در Redis ذخیره می‌کنیم تا لینک و تعداد را دریافت کنیم
}
