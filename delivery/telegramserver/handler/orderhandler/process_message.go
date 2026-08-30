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
	"telegram-service-platform/pkg/richerror"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// handleMessage پیام‌های متنی کاربر را در حین فرآیند سفارش پردازش می‌کند
func (h *Handler) handleMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "orderhandler.handleMessage"

	if update.Message == nil || update.Message.Text == "" {
		return
	}

	chatID := update.Message.Chat.ID
	telegramID := update.Message.From.ID
	text := strings.TrimSpace(update.Message.Text)

	stateResp, err := h.orderFlowService.GetOrderFlow(ctx, orderparams.GetOrderFlowRequest{TelegramID: entity.TelegramId(telegramID)})
	if err != nil || stateResp == nil {

		if !strings.HasPrefix(text, "/") {
			_ = h.messenger.Send(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "⚠️ لطفاً ابتدا از منوی اصلی یک سرویس را انتخاب کنید.",
			})
		}
		return
	}

	state := stateResp.Stage

	switch state {
	case orderentity.OrderFlowStageWaitingForLink:
		h.handleLinkInput(ctx, chatID, telegramID, text, stateResp, op)
	case orderentity.OrderFlowStageWaitingForQuantity:
		h.handleQuantityInput(ctx, chatID, telegramID, text, stateResp, op)
	default:

		break
	}
}

func (h *Handler) handleLinkInput(ctx context.Context, chatID, telegramID int64, text string, state *orderentity.OrderFlowState, op string) {

	req := orderparams.SubmitLinkRequest{
		Link: text,
	}

	if err := h.validator.ValidateLink(req); err != nil {
		if richErr, ok := errors.AsType[*richerror.RichError](err); ok {
			_ = h.messenger.Send(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "❌ " + richErr.Message(),
			})
			logger.Logger.Warn("link validation failed",
				zap.String("op", op),
				zap.Any("meta", richErr.Meta()),
			)
		}
		return
	}

	state.Link = text
	state.Stage = orderentity.OrderFlowStageConfirming

	if err := h.orderFlowService.SaveOrderFlow(ctx, orderparams.SaveOrderFlowRequest{
		TelegramID: entity.TelegramId(telegramID),
		State:      *state,
		TTLMins:    10,
	}); err != nil {
		logger.Logger.Error("failed to save order flow state (link)", zap.String("op", op), zap.Error(err))
		h.handleError(ctx, chatID, op, err)
		return
	}

	h.showConfirmOrder(ctx, chatID, state)
}

func (h *Handler) handleQuantityInput(ctx context.Context, chatID, telegramID int64, text string, state *orderentity.OrderFlowState, op string) {
	quantity, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		_ = h.messenger.Send(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "❌ تعداد نامعتبر است. لطفاً فقط عدد وارد کنید.",
		})
		return
	}

	req := orderparams.SubmitQuantityRequest{
		Quantity: quantity,
		Min:      state.MinQuantity,
		Max:      state.MaxQuantity,
	}

	if err := h.validator.ValidateQuantity(req); err != nil {
		if richErr, ok := errors.AsType[*richerror.RichError](err); ok {

			_ = h.messenger.Send(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "❌ " + richErr.Message(),
			})

			logger.Logger.Warn("quantity validation failed",
				zap.String("op", op),
				zap.Any("meta", richErr.Meta()),
			)
		}
		return
	}

	// ۱. محاسبه قیمت دلاری: (Rate * Quantity) / 1000
	usdPrice := state.Rate.Mul(entity.Amount(decimal.NewFromInt(quantity))).Div(entity.Amount(decimal.NewFromInt(1000)))

	// ۲. دریافت نرخ دلار به تومان از PriceService
	tomanRate, gErr := h.pricingSvc.GetUsdTomanPrice(ctx)

	if gErr != nil {
		logger.Logger.Error("failed to get usd to toman price", zap.String("op", op), zap.Error(gErr))
		_ = h.messenger.Send(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "❌ خطا در دریافت نرخ لحظه‌ای ارز. لطفاً چند لحظه بعد دوباره تلاش کنید.",
		})
		return
	}

	// ۳. محاسبه قیمت نهایی تومانی و گرد کردن به ۲ رقم اعشار (یا رند کردن به عدد صحیح)
	tomanPrice := usdPrice.Mul(tomanRate).Round(3)

	// ۴. به‌روزرسانی State
	state.Quantity = quantity
	state.Price = tomanPrice
	state.Stage = orderentity.OrderFlowStageWaitingForLink

	if err := h.orderFlowService.SaveOrderFlow(ctx, orderparams.SaveOrderFlowRequest{
		TelegramID: entity.TelegramId(telegramID),
		State:      *state,
		TTLMins:    10,
	}); err != nil {
		logger.Logger.Error("failed to save order flow state (quantity)", zap.String("op", op), zap.Error(err))
		h.handleError(ctx, chatID, op, err)
		return
	}

	tomanPriceStr := tomanPrice.String()
	usdPriceStr := usdPrice.String()

	message := fmt.Sprintf(
		"✅ تعداد <b>%d</b> با موفقیت ثبت شد.\n\n"+
			"💰 <b>قیمت برآورد شده برای این سفارش:</b>\n"+
			"• به دلار: <code>%s $</code>\n"+
			"• به تومان: <code>%s تومان</code>\n\n"+
			"🔗 حالا لطفاً لینک کانال یا گروه خود را ارسال کنید.\n"+
			"(مثال: <code>https://t.me/YourChannel</code>)\n\n"+
			"⚠️ <b>توجه:</b> لینک باید عمومی (Public) باشد تا سرویس قابل انجام باشد.",
		quantity,
		usdPriceStr,
		tomanPriceStr,
	)

	sErr := h.messenger.Send(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      message,
		ParseMode: models.ParseModeHTML,
	})

	fmt.Println("ERRORS : ", sErr.Error())
}

func (h *Handler) showConfirmOrder(ctx context.Context, chatID int64, state *orderentity.OrderFlowState) {
	message := fmt.Sprintf(
		"📋 خلاصه سفارش شما:\n\n"+
			"📱 پلتفرم: %s\n"+
			"📦 سرویس: %s\n"+
			"🔢 تعداد: %d عدد\n"+
			"🔗 لینک: %s\n\n"+
			"💰 مبلغ قابل پرداخت: %s تومان\n\n"+
			"آیا اطلاعات فوق صحیح است؟",
		state.Platform,
		state.ServiceName,
		state.Quantity,
		state.Link,
		state.Price.String(), // نمایش قیمت با فرمت decimal
	)

	_ = h.messenger.Send(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        message,
		ReplyMarkup: keyboard.OrderConfirmMenu(),
	})
}

func (h *Handler) handleError(ctx context.Context, chatID int64, op string, err error) {
	_ = h.messenger.Send(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "❌ خطایی در پردازش درخواست شما رخ داد. لطفاً دوباره تلاش کنید.",
	})
	logger.Logger.Error("orderhandler error", zap.String("op", op), zap.Error(err))
}
