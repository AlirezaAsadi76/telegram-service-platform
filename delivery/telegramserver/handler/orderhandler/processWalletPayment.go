package orderhandler

import (
	"context"
	"errors"
	"fmt"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/params/orderparams"
	"telegram-service-platform/params/userparams"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"

	"telegram-service-platform/entity"
	"telegram-service-platform/entity/productentity"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/checkoutparams"
	"telegram-service-platform/pkg/richerror"
)

func (h *Handler) processWalletPayment(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "orderhandler.processWalletPayment"

	if update.CallbackQuery == nil || update.CallbackQuery.Message.Message == nil {
		return
	}

	chatID := update.CallbackQuery.Message.Message.Chat.ID
	telegramID := update.CallbackQuery.From.ID

	stateResp, err := h.orderFlowService.GetOrderFlow(ctx, orderparams.GetOrderFlowRequest{TelegramID: entity.TelegramId(telegramID)})
	if err != nil || stateResp == nil || stateResp.Stage != orderentity.OrderFlowStageConfirming {
		_ = h.messenger.Send(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "⚠️ سفارش معتبری برای پرداخت یافت نشد. لطفاً از ابتدا شروع کنید.",
		})
		return
	}

	state := stateResp
	user, uErr := h.userService.FindUserByTelegramID(ctx, userparams.FindUserByTelegramIDRequest{TelegramID: entity.TelegramId(telegramID)})
	if uErr != nil || !user.Found {
		_ = h.messenger.Send(ctx, &bot.SendMessageParams{})
	}

	if err := h.ensurePriceLock(
		ctx,
		telegramID,
		user.UserInfo.Id,
		state,
		entity.PriceLockPaymentMethodWallet,
	); err != nil {
		if errors.Is(err, ErrPriceLockMethodMismatch) {
			_ = h.messenger.Send(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   priceLockMethodMismatchMessage,
			})
			return
		}

		logger.Logger.Error(
			"failed to lock Product price for wallet payment",
			zap.String("op", op),
			zap.Uint64("user_id", user.UserInfo.Id),
			zap.Uint64("mapping_id", state.ServiceID),
			zap.Error(err),
		)

		_ = h.messenger.Send(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   priceLockFailedMessage,
		})
		return
	}

	req := checkoutparams.WalletPurchaseRequest{
		UserID:                 user.UserInfo.Id,
		ProductType:            productentity.ProductTypeSMM,
		ProductID:              state.ServiceID,
		Quantity:               state.Quantity,
		TargetLink:             state.Link,
		Amount:                 state.Price,
		Currency:               state.Currency,
		IdempotencyKey:         state.PurchaseID,
		PriceLockedAt:          state.PriceLockedAt,
		PriceLockExpiresAt:     state.PriceLockExpiresAt,
		PriceLockPaymentMethod: entity.PriceLockPaymentMethodWallet,
	}

	if vErr := h.validator.ValidatePaymentPriceLock(req, entity.PriceLockPaymentMethodWallet); vErr != nil {
		_ = h.messenger.Send(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   vErr.Error(),
		})
	}
	_, chErr := h.checkoutService.ProcessWalletPurchase(
		ctx,
		checkoutparams.WalletPurchaseRequest{
			UserID:                 user.UserInfo.Id,
			ProductType:            productentity.ProductTypeSMM,
			ProductID:              state.ServiceID,
			Quantity:               state.Quantity,
			TargetLink:             state.Link,
			Amount:                 state.Price,
			Currency:               state.Currency,
			IdempotencyKey:         state.PurchaseID,
			PriceLockedAt:          state.PriceLockedAt,
			PriceLockExpiresAt:     state.PriceLockExpiresAt,
			PriceLockPaymentMethod: entity.PriceLockPaymentMethodWallet,
		},
	)

	if chErr != nil {

		var rErr *richerror.RichError
		isRichError := errors.As(chErr, &rErr)

		if isRichError && rErr.Kind() == richerror.KindValidation {
			_ = h.messenger.Send(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "❌ موجودی کیف پول شما کافی نیست.\nلطفاً ابتدا کیف پول خود را شارژ کنید.\n\n(سفارش شما تا ۱۰ دقیقه دیگر در سیستم باقی می‌ماند تا پس از شارژ، پرداخت را انجام دهید.)",
			})

		} else {
			logger.Logger.Error("wallet purchase failed",
				zap.String("op", op),
				zap.Int64("telegram_id", telegramID),
				zap.Error(chErr),
			)
			_ = h.messenger.Send(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "❌ خطایی در پردازش پرداخت رخ داد. لطفاً دوباره تلاش کنید یا با پشتیبانی تماس بگیرید.",
			})
		}
		return
	}

	_ = h.messenger.Send(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text: fmt.Sprintf(
			"✅ <b>پرداخت با موفقیت انجام شد!</b>\n\n"+
				"🎉 سفارش شما ثبت گردید و در حال پردازش است.\n\n"+
				"💰 مبلغ کسر شده: <code>%s</code> تومان\n\n"+
				"📌 می‌توانید وضعیت سفارش خود را از بخش «💳 تراکنش‌ها» پیگیری کنید.",
			state.Price.String(),
		),
		ParseMode: models.ParseModeHTML,
	})

	logger.Logger.Info("wallet purchase completed successfully",
		zap.String("op", op),
		zap.Int64("telegram_id", telegramID),
		zap.String("price", state.Price.String()),
	)
}
