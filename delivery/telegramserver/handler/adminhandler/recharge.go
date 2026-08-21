package adminhandler

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"

	"telegram-service-platform/entity"
	"telegram-service-platform/logger"
	"telegram-service-platform/params"
	"telegram-service-platform/params/checkoutparams"
)

func (h *Handler) Recharge(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "adminhandler.Recharge"

	if update.Message == nil || update.Message.Text == "" {
		return
	}

	chatID := update.Message.Chat.ID
	adminTelegramID := update.Message.From.ID

	parts := strings.Fields(update.Message.Text)
	if len(parts) != 3 {
		h.messenger.Send(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "⚠️ فرمت اشتباه.\nاستفاده صحیح: `/recharge <telegram_id> <amount>`\nمثال: `/recharge 123456789 100000`",
		})
		return
	}

	telegramID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		h.messenger.Send(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "⚠️ Telegram ID نامعتبر. باید یک عدد باشد.",
		})
		return
	}

	amountValue, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil || amountValue <= 0 {
		h.messenger.Send(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "⚠️ مبلغ نامعتبر. باید یک عدد مثبت باشد.",
		})
		return
	}

	userResp, userErr := h.userService.FindUserByTelegramID(ctx, params.FindUserByTelegramIDRequest{
		TelegramID: entity.TelegramId(telegramID),
	})
	if userErr != nil {
		logger.Logger.Error("find user failed",
			zap.String("op", op),
			zap.Int64("telegram_id", telegramID),
			zap.Error(userErr),
		)
		h.messenger.Send(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "❌ خطا در جستجوی کاربر.",
		})
		return
	}

	if !userResp.Found {
		h.messenger.Send(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   fmt.Sprintf("⚠️ کاربر با Telegram ID %d یافت نشد.", telegramID),
		})
		return
	}

	amount := entity.Amount(amountValue)
	rechargeErr := h.checkoutService.ProcessManualWalletRecharge(ctx, checkoutparams.ManualRechargeRequest{
		AdminID:  uint64(adminTelegramID),
		UserID:   userResp.UserInfo.Id,
		Amount:   amount,
		Currency: entity.CurrencyTOMAN,
	})

	if rechargeErr != nil {
		logger.Logger.Error("recharge failed",
			zap.String("op", op),
			zap.Uint64("user_id", userResp.UserInfo.Id),
			zap.Int64("amount", amountValue),
			zap.Error(rechargeErr),
		)
		h.messenger.Send(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "❌ خطا در شارژ کیف پول.",
		})
		return
	}

	successMsg := fmt.Sprintf(
		"✅ کیف پول با موفقیت شارژ شد.\n\n"+
			"👤 Telegram ID: %d\n"+
			"👤 Username: @%s\n"+
			"💰 مبلغ: %d تومان",
		telegramID,
		userResp.UserInfo.Username,
		amountValue,
	)

	h.messenger.Send(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   successMsg,
	})
}
