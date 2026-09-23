package orderhandler

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"telegram-service-platform/delivery/telegramserver/keyboard"
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/orderparams"
	"telegram-service-platform/params/productparams"
	"telegram-service-platform/pkg/richerror"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

func (h *Handler) handleMessage(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	const op = "orderhandler.handleMessage"

	if update.Message == nil ||
		update.Message.Text == "" {
		return
	}

	chatID := update.Message.Chat.ID
	telegramID := update.Message.From.ID
	text := strings.TrimSpace(update.Message.Text)

	stateResp, err := h.orderFlowService.GetOrderFlow(
		ctx,
		orderparams.GetOrderFlowRequest{
			TelegramID: entity.TelegramId(telegramID),
		},
	)
	if err != nil || stateResp == nil {
		if !strings.HasPrefix(text, "/") {
			_ = h.messenger.Send(
				ctx,
				&bot.SendMessageParams{
					ChatID: chatID,
					Text:   "⚠️ لطفاً ابتدا از منوی اصلی یک سرویس را انتخاب کنید.",
				},
			)
		}

		return
	}

	switch stateResp.Stage {
	case orderentity.OrderFlowStageWaitingForLink:
		h.handleLinkInput(
			ctx,
			chatID,
			telegramID,
			text,
			stateResp,
			op,
		)

	case orderentity.OrderFlowStageWaitingForQuantity:
		h.handleQuantityInput(
			ctx,
			chatID,
			telegramID,
			text,
			stateResp,
			op,
		)
	}
}

func (h *Handler) handleLinkInput(
	ctx context.Context,
	chatID int64,
	telegramID int64,
	text string,
	state *orderentity.OrderFlowState,
	op string,
) {
	req := orderparams.SubmitLinkRequest{
		Link: text,
	}

	if err := h.validator.ValidateLink(req); err != nil {
		if richErr, ok := errors.AsType[*richerror.RichError](err); ok {
			_ = h.messenger.Send(
				ctx,
				&bot.SendMessageParams{
					ChatID: chatID,
					Text:   "❌ " + richErr.Message(),
				},
			)

			logger.Logger.Warn(
				"link validation failed",
				zap.String("op", op),
				zap.Any("meta", richErr.Meta()),
			)
		}

		return
	}

	state.Link = text
	state.Stage = orderentity.OrderFlowStageConfirming

	if err := h.orderFlowService.SaveOrderFlow(
		ctx,
		orderparams.SaveOrderFlowRequest{
			TelegramID: entity.TelegramId(telegramID),
			State:      *state,
			TTLMins:    10,
		},
	); err != nil {
		logger.Logger.Error(
			"failed to save order flow state (link)",
			zap.String("op", op),
			zap.Error(err),
		)

		h.handleError(ctx, chatID, op, err)
		return
	}

	h.showConfirmOrder(
		ctx,
		chatID,
		state,
	)
}

func (h *Handler) handleQuantityInput(
	ctx context.Context, chatID int64, telegramID int64,
	text string, state *orderentity.OrderFlowState, op string) {
	quantity, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		_ = h.messenger.Send(
			ctx,
			&bot.SendMessageParams{
				ChatID: chatID,
				Text:   "❌ تعداد نامعتبر است. لطفاً فقط عدد وارد کنید.",
			},
		)

		return
	}

	validationReq := orderparams.SubmitQuantityRequest{
		Quantity: quantity,
		Min:      state.MinQuantity,
		Max:      state.MaxQuantity,
	}

	if err := h.validator.ValidateQuantity(validationReq); err != nil {
		if richErr, ok := errors.AsType[*richerror.RichError](err); ok {
			_ = h.messenger.Send(
				ctx,
				&bot.SendMessageParams{
					ChatID: chatID,
					Text:   "❌ " + richErr.Message(),
				},
			)

			logger.Logger.Warn(
				"quantity validation failed",
				zap.String("op", op),
				zap.Any("meta", richErr.Meta()),
			)
		}

		return
	}

	priceResponse, cErr := h.productService.CalculateSMMPrice(
		ctx,
		productparams.CalculateSMMPriceRequest{
			MappingID: int64(state.ServiceID),
			Quantity:  quantity,
		},
	)
	if cErr != nil {
		logger.Logger.Error(
			"failed to calculate SMM price",
			zap.String("op", op),
			zap.Uint64("mapping_id", state.ServiceID),
			zap.Int64("quantity", quantity),
			zap.Error(cErr),
		)

		_ = h.messenger.Send(
			ctx,
			&bot.SendMessageParams{
				ChatID: chatID,
				Text:   "❌ خطایی در محاسبه قیمت سفارش رخ داد. لطفاً چند لحظه بعد دوباره تلاش کنید.",
			},
		)

		return
	}

	state.Quantity = quantity
	state.Rate = priceResponse.Rate
	state.Price = priceResponse.Price.Toman
	state.Currency = entity.CurrencyTOMAN
	state.Stage = orderentity.OrderFlowStageWaitingForLink

	if err := h.orderFlowService.SaveOrderFlow(
		ctx,
		orderparams.SaveOrderFlowRequest{
			TelegramID: entity.TelegramId(telegramID),
			State:      *state,
			TTLMins:    10,
		},
	); err != nil {
		logger.Logger.Error(
			"failed to save order flow state (quantity)",
			zap.String("op", op),
			zap.Error(err),
		)

		h.handleError(ctx, chatID, op, err)

		return
	}

	message := fmt.Sprintf(
		"✅ تعداد <b>%d</b> با موفقیت ثبت شد.\n\n"+
			"💰 <b>قیمت برآورد شده برای این سفارش:</b>\n"+
			"• به دلار: <code>%s $</code>\n"+
			"• به تومان: <code>%s تومان</code>\n\n"+
			"🔗 حالا لطفاً لینک کانال یا گروه خود را ارسال کنید.\n"+
			"(مثال: <code>https://t.me/YourChannel</code>)\n\n"+
			"⚠️ <b>توجه:</b> لینک باید عمومی (Public) باشد تا سرویس قابل انجام باشد.",
		quantity,
		priceResponse.Price.USD.String(),
		priceResponse.Price.Toman.String(),
	)

	if err := h.messenger.Send(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      message,
		ParseMode: models.ParseModeHTML,
	},
	); err != nil {
		logger.Logger.Warn(
			"failed to send quantity result message",
			zap.String("op", op),
			zap.Int64("telegram_id", telegramID),
			zap.Error(err),
		)
	}
}

func (h *Handler) showConfirmOrder(ctx context.Context, chatID int64, state *orderentity.OrderFlowState) {
	message := fmt.Sprintf(
		"📋 <b>خلاصه سفارش شما:</b>\n\n"+
			"📱 پلتفرم: %s\n"+
			"📦 سرویس: %s\n"+
			"🔢 تعداد: <code>%d</code> عدد\n"+
			"🔗 لینک: %s\n\n"+
			"💰 <b>مبلغ قابل پرداخت:</b> <code>%s تومان</code>\n\n"+
			"آیا اطلاعات فوق صحیح است؟",
		state.Platform,
		state.ServiceName,
		state.Quantity,
		state.Link,
		state.Price.String(),
	)

	_ = h.messenger.Send(
		ctx,
		&bot.SendMessageParams{
			ChatID:      chatID,
			Text:        message,
			ReplyMarkup: keyboard.OrderConfirmMenu(),
			ParseMode:   models.ParseModeHTML,
		},
	)
}

func (h *Handler) handleError(ctx context.Context, chatID int64, op string, err error) {
	_ = h.messenger.Send(
		ctx,
		&bot.SendMessageParams{
			ChatID: chatID,
			Text:   "❌ خطایی در پردازش درخواست شما رخ داد. لطفاً دوباره تلاش کنید.",
		},
	)

	logger.Logger.Error(
		"orderHandler error",
		zap.String("op", op),
		zap.Error(err),
	)
}
