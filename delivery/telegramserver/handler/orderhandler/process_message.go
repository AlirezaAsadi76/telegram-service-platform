package orderhandler

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"telegram-service-platform/delivery/telegramserver/keyboard"
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/orderparams"

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
	if !strings.HasPrefix(text, "http://") && !strings.HasPrefix(text, "https://") {
		_ = h.messenger.Send(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "❌ لینک نامعتبر است.\nلطفاً لینک را با http:// یا https:// شروع کنید.",
		})
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
	if err != nil || quantity <= 0 {
		_ = h.messenger.Send(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "❌ تعداد نامعتبر است.\nلطفاً فقط یک عدد مثبت وارد کنید.",
		})
		return
	}

	if quantity < state.MinQuantity || quantity > state.MaxQuantity {
		_ = h.messenger.Send(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   fmt.Sprintf("❌ تعداد خارج از محدوده مجاز است.\n📊 حداقل: %d\n📊 حداکثر: %d", state.MinQuantity, state.MaxQuantity),
		})
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
	tomanPrice := usdPrice.Mul(tomanRate).Round(2)

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

	// ۵. پیام درخواست لینک (UX بهبود یافته)
	_ = h.messenger.Send(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text: fmt.Sprintf(
			"✅ تعداد «%d» ثبت شد.\n\n"+
				"🔗 حالا لطفاً لینک کانال یا گروه خود را ارسال کنید.\n"+
				"(مثال: https://t.me/YourChannel)\n\n"+
				"⚠️ توجه: لینک باید عمومی (Public) باشد تا سرویس قابل انجام باشد.",
			quantity,
		),
	})
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
