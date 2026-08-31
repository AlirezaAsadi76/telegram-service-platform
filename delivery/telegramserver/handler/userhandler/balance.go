package userhandler

import (
	"context"
	"errors"
	"fmt"
	"telegram-service-platform/entity"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"

	"telegram-service-platform/logger"
	"telegram-service-platform/params/walletparam"
	"telegram-service-platform/pkg/richerror"
)

func (h Handler) WalletBalance(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "userhandler.handleWalletBalance"

	telegramID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	userID, err := h.userValidator.ValidationUserExistence(ctx, entity.TelegramId(telegramID))
	if err != nil {
		if richErr, ok := errors.AsType[*richerror.RichError](err); ok {

			_ = h.messenger.Send(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "❌ " + "خطایی در پردازش رخ داده است",
			})

			logger.Logger.Warn("balance validation failed",
				zap.String("op", op),
				zap.Int64("telegram_id", telegramID),
				zap.Any("meta", richErr.Meta()),
			)
		}
		return
	}

	balanceResp, gErr := h.walletService.GetBalance(ctx, walletparam.GetBalanceRequest{UserID: userID})
	if gErr != nil {
		logger.Logger.Error("failed to get wallet balance",
			zap.String("op", op),
			zap.Uint64("user_id", userID),
			zap.Error(gErr),
		)
		_ = h.messenger.Send(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "❌ خطایی در دریافت موجودی رخ داد. لطفاً دوباره تلاش کنید.",
		})
		return
	}

	// ۳. نمایش موجودی با فرمت خوانا
	message := fmt.Sprintf(
		"💰 <b>کیف پول شما</b>\n\n"+
			"🏦 موجودی فعلی: <code>%s</code> %s\n\n"+
			"برای شارژ کیف پول، لطفاً از منوی اصلی اقدام کنید یا با پشتیبانی در ارتباط باشید.",
		balanceResp.Balance.String(),
		balanceResp.Currency,
	)

	_ = h.messenger.Send(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      message,
		ParseMode: models.ParseModeHTML,
	})
}
