package orderhandler

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"telegram-service-platform/delivery/telegramserver/keyboard"
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/params/orderparams"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
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
		h.handleLinkInput(ctx, chatID, telegramID, text, stateResp)
	case orderentity.OrderFlowStageWaitingForQuantity:
		h.handleQuantityInput(ctx, chatID, telegramID, text, stateResp)
	default:

		break
	}
}

func (h *Handler) handleLinkInput(ctx context.Context, chatID, telegramID int64, text string, state *orderentity.OrderFlowState) {
	if !strings.HasPrefix(text, "http://") && !strings.HasPrefix(text, "https://") {
		_ = h.messenger.Send(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "❌ لینک نامعتبر است.\nلطفاً لینک را با http:// یا https:// شروع کنید.",
		})
		return
	}

	state.Link = text
	state.Stage = orderentity.OrderFlowStageWaitingForQuantity

	_ = h.orderFlowService.SaveOrderFlow(ctx, orderparams.SaveOrderFlowRequest{
		TelegramID: entity.TelegramId(telegramID),
		State:      *state,
		TTLMins:    10,
	})

	_ = h.messenger.Send(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   fmt.Sprintf("✅ لینک ثبت شد.\n\nحالا تعداد مورد نظر را وارد کنید:\n(حداقل: %d | حداکثر: %d)", state.MinQuantity, state.MaxQuantity),
	})
}

func (h *Handler) handleQuantityInput(ctx context.Context, chatID, telegramID int64, text string, state *orderentity.OrderFlowState) {
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
			Text:   fmt.Sprintf("❌ تعداد خارج از محدوده مجاز است.\nحداقل: %d | حداکثر: %d", state.MinQuantity, state.MaxQuantity),
		})
		return
	}

	// محاسبه قیمت (همه داده‌ها از Redis خوانده شده و نیازی به DB نیست!)
	price := (state.Rate * entity.Amount(quantity)) / 1000
	state.Quantity = quantity
	state.Price = price
	state.Stage = orderentity.OrderFlowStageConfirming

	_ = h.orderFlowService.SaveOrderFlow(ctx, orderparams.SaveOrderFlowRequest{
		TelegramID: entity.TelegramId(telegramID),
		State:      *state,
		TTLMins:    10,
	})

	h.showConfirmOrder(ctx, chatID, state)
}

func (h *Handler) showConfirmOrder(ctx context.Context, chatID int64, state *orderentity.OrderFlowState) {
	message := fmt.Sprintf(
		"✅ خلاصه سفارش شما:\n\n"+
			"📱 پلتفرم: %s\n📂 دسته: %s\n📦 سرویس: %s\n\n"+
			"🔗 لینک: %s\n🔢 تعداد: %d\n💰 قیمت کل: %d تومان\n\n"+
			"لطفاً روش پرداخت را انتخاب کنید:",
		state.Platform, state.Category, state.ServiceName, state.Link, state.Quantity, state.Price,
	)

	_ = h.messenger.Send(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        message,
		ReplyMarkup: keyboard.OrderConfirmMenu(),
	})
}
