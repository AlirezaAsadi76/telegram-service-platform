package mainhandler

import (
	"context"
	"errors"
	"strconv"
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/params/orderparams"
	"telegram-service-platform/pkg/helpers"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/google/uuid"
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
	telegramID := update.CallbackQuery.From.ID

	h.clearActiveOrderFlowIfAny(ctx, telegramID, op)

	data, ccErr := h.callbackQueryData(update.CallbackQuery.Data, CategorySplitMode)
	platformName := data[platformSplitMode]
	categoryName := data[CategorySplitMode]
	serviceIDStr := data[ServiceSplitMode]
	if ccErr != nil {
		h.handleError(ctx, chatID, op, errors.New("invalid callback data"))
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
	if gErr != nil || mappingResp == nil || mappingResp.SmmMapping == nil {
		logger.Logger.Error("failed to get service mapping details",
			zap.String("op", op),
			zap.Int64("serviceMappingID", serviceID),
			zap.Error(gErr),
		)
		if gErr != nil {
			h.handleError(ctx, chatID, op, gErr)
		} else {
			h.handleError(ctx, chatID, op, errors.New("service mapping not found"))
		}
		return
	}

	serviceResp, gsErr := h.productService.GetSMMServiceByID(ctx, productparams.GetSmmServiceByIDRequest{
		Id: mappingResp.SmmMapping.SmmServiceId,
	})
	if gsErr != nil || serviceResp == nil || serviceResp.Smm == nil {
		logger.Logger.Error("failed to get service details",
			zap.String("op", op),
			zap.Int64("serviceID", serviceID),
			zap.Error(gsErr),
		)
		if gsErr != nil {
			h.handleError(ctx, chatID, op, gsErr)
		} else {
			h.handleError(ctx, chatID, op, errors.New("smm service not found"))
		}
		return
	}

	state := orderentity.OrderFlowState{
		Stage:      orderentity.OrderFlowStageWaitingForQuantity,
		PurchaseID: uuid.NewString(),

		Platform:    platformName,
		Category:    categoryName,
		ServiceID:   uint64(serviceID),
		ServiceName: mappingResp.SmmMapping.ButtonName,
		MinQuantity: serviceResp.Smm.Min,
		MaxQuantity: serviceResp.Smm.Max,
		Rate:        serviceResp.Smm.Rate,
		Currency:    entity.CurrencyTOMAN,
	}

	saveErr := h.orderFlowService.SaveOrderFlow(ctx, orderparams.SaveOrderFlowRequest{
		TelegramID: entity.TelegramId(chatID),
		State:      state,
	})
	if saveErr != nil {
		logger.Logger.Error("failed to save order flow state",
			zap.String("op", op),
			zap.Int64("telegram_id", chatID),
			zap.Error(saveErr),
		)
		h.handleError(ctx, chatID, op, saveErr)
		return
	}

	message := serviceSelectedMessage(
		mappingResp.SmmMapping.ButtonName,
		helpers.GetPlatformDisplayName(platformName),
		helpers.GetCategoryIcon(categoryName),
		helpers.GetCategoryDisplayName(categoryName),
		serviceResp.Smm.Min,
		serviceResp.Smm.Max,
	)

	if sendErr := h.messenger.Edit(ctx, &bot.EditMessageTextParams{
		ChatID:    chatID,
		Text:      message,
		MessageID: messageID,
	}); sendErr != nil {
		logger.Logger.Error("failed to send service details", zap.Error(sendErr), zap.String("op", op))
	}

}
