package orderhandler

import (
	"context"
	"fmt"
	"telegram-service-platform/delivery/telegramserver/keyboard"

	"telegram-service-platform/entity"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/entity/productentity"
	"telegram-service-platform/logger"
	"telegram-service-platform/params"
	"telegram-service-platform/params/checkoutparams"
	"telegram-service-platform/params/orderparams"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

func (h *Handler) processDirectPayment(ctx context.Context, b *bot.Bot, update *models.Update, method paymententity.PaymentMethod, methodName string) {
	const op = "orderhandler.processDirectPayment"

	if update.CallbackQuery == nil || update.CallbackQuery.Message.Message == nil {
		return
	}

	chatID := update.CallbackQuery.Message.Message.Chat.ID
	telegramID := update.CallbackQuery.From.ID

	// ۱. دریافت State
	stateResp, gErr := h.orderFlowService.GetOrderFlow(ctx, orderparams.GetOrderFlowRequest{
		TelegramID: entity.TelegramId(telegramID),
	})

	if gErr != nil || stateResp == nil || stateResp.Stage != orderentity.OrderFlowStageConfirming {
		_ = h.messenger.Send(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "⚠️ سفارش معتبری برای پرداخت یافت نشد. لطفاً از ابتدا شروع کنید.",
		})
		return
	}

	state := stateResp
	user, uErr := h.userService.FindUserByTelegramID(ctx, params.FindUserByTelegramIDRequest{TelegramID: entity.TelegramId(telegramID)})
	if uErr != nil || !user.Found {
		_ = h.messenger.Send(ctx, &bot.SendMessageParams{})
	}

	resp, pErr := h.checkoutService.ProcessDirectPaymentPurchase(ctx, checkoutparams.DirectPaymentPurchase{
		UserID:      user.UserInfo.Id,
		ProductType: productentity.SMM,
		ProductID:   state.ServiceID,
		Quantity:    state.Quantity,
		TargetLink:  state.Link,
		Amount:      state.Price,
		Currency:    state.Currency,
		Method:      method,
	})

	if pErr != nil {
		logger.Logger.Error("direct payment purchase failed",
			zap.String("op", op),
			zap.Int64("telegram_id", telegramID),
			zap.String("method", string(method)),
			zap.Error(pErr),
		)
		_ = h.messenger.Send(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "❌ خطایی در ایجاد درخواست پرداخت رخ داد. لطفاً دوباره تلاش کنید.",
		})
		return
	}

	message := fmt.Sprintf(
		"💳 <b>درخواست %s ثبت شد</b>\n\n"+
			"📋 شماره سفارش: <code>%d</code>\n"+
			"💰 مبلغ: <code>%s</code> تومان\n\n"+
			"برای تکمیل پرداخت، روی دکمه زیر کلیک کنید:\n\n"+
			"<i>⚠️ توجه: در حال حاضر این بخش در محیط شبیه‌سازی (Stub) قرار دارد.</i>",
		methodName,
		resp.OrderID,
		state.Price.String(),
	)

	_ = h.messenger.Send(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        message,
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: keyboard.PaymentActionMenu(resp.PaymentURL),
	})

	logger.Logger.Info("direct payment url generated",
		zap.String("op", op),
		zap.Int64("telegram_id", telegramID),
		zap.Uint64("order_id", resp.OrderID),
		zap.Uint64("payment_id", resp.PaymentID),
		zap.String("method", string(method)),
	)
}

func (h *Handler) checkPaymentStatus(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "orderhandler.checkPaymentStatus"

	if update.CallbackQuery == nil || update.CallbackQuery.Message.Message == nil {
		return
	}

	chatID := update.CallbackQuery.Message.Message.Chat.ID

	_ = h.messenger.Send(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text: "⏳ <b>وضعیت پرداخت:</b>\n\n" +
			"لطفاً چند لحظه صبر کنید تا سیستم پرداخت را بررسی کند.\n\n" +
			"<i>(در نسخه نهایی، این دکمه مستقیماً وضعیت را از OrderService استعلام می‌گیرد)</i>",
		ParseMode: models.ParseModeHTML,
	})
}
