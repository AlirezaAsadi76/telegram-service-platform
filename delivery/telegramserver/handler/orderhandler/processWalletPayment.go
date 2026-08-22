package orderhandler

import (
	"context"
	"errors"
	"fmt"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/params/orderparams"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"

	"telegram-service-platform/entity"
	"telegram-service-platform/entity/productentity"
	"telegram-service-platform/logger"
	"telegram-service-platform/params"
	"telegram-service-platform/params/checkoutparams"
	"telegram-service-platform/pkg/richerror"
)

func (h *Handler) processWalletPayment(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "orderhandler.processWalletPayment"

	if update.CallbackQuery == nil || update.CallbackQuery.Message.Message == nil {
		return
	}

	chatID := update.CallbackQuery.Message.Message.ID
	telegramID := update.CallbackQuery.From.ID

	// ۱. دریافت State نهایی
	stateResp, err := h.orderFlowService.GetOrderFlow(ctx, orderparams.GetOrderFlowRequest{TelegramID: entity.TelegramId(telegramID)})
	if err != nil || stateResp == nil || stateResp.Stage != orderentity.OrderFlowStageConfirming {
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
	// ۲. فراخوانی CheckoutService
	err = h.checkoutService.ProcessWalletPurchase(ctx, checkoutparams.WalletPurchaseRequest{
		UserID:      user.UserInfo.Id,
		ProductType: productentity.SMM,
		ProductID:   state.ServiceID,
		Quantity:    state.Quantity,
		TargetLink:  state.Link,
		Amount:      state.Price,
		Currency:    state.Currency,
	})

	if err != nil {
		// بررسی خطای کمبود موجودی برای پیام کاربرپسندتر
		if richErr, ok := errors.AsType[*richerror.RichError](err); ok && richErr.Kind() == richerror.KindValidation {
			_ = h.messenger.Send(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "❌ موجودی کیف پول شما کافی نیست.\nلطفاً ابتدا کیف پول خود را شارژ کنید.\n\n(سفارش شما تا ۱۰ دقیقه دیگر در سیستم باقی می‌ماند تا پس از شارژ، پرداخت را انجام دهید.)",
			})
			// State را پاک نمی‌کنیم تا کاربر بعد از شارژ بتواند برگردد و پرداخت کند
		} else {
			logger.Logger.Error("wallet purchase failed",
				zap.String("op", op),
				zap.Int64("telegram_id", telegramID),
				zap.Error(err),
			)
			_ = h.messenger.Send(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "❌ خطایی در پردازش پرداخت رخ داد. لطفاً دوباره تلاش کنید یا با پشتیبانی تماس بگیرید.",
			})
		}
		return
	}

	_ = h.orderFlowService.CompleteOrderFlow(ctx, orderparams.DeleteOrderFlowRequest{
		TelegramID: entity.TelegramId(telegramID),
	}, 0)

	_ = h.messenger.Send(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   fmt.Sprintf("✅ پرداخت با موفقیت انجام شد!\n\nسفارش شما ثبت گردید و در حال پردازش است.\n💰 مبلغ کسر شده: %d تومان", state.Price),
	})
}
