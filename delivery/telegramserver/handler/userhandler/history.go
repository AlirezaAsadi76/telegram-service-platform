package userhandler

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"telegram-service-platform/entity"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"

	"telegram-service-platform/logger"
	"telegram-service-platform/params/orderparams"
	"telegram-service-platform/pkg/richerror"
)

func (h Handler) TransactionsHistory(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "userhandler.TransactionsHistory"

	if update.Message == nil {
		return
	}

	telegramID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	userID, err := h.userValidator.ValidationUserExistence(ctx, entity.TelegramId(telegramID))
	if err != nil {
		if richErr, ok := errors.AsType[*richerror.RichError](err); ok {
			_ = h.messenger.Send(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "❌ " + richErr.Message(),
			})
			logger.Logger.Warn("history validation failed",
				zap.String("op", op),
				zap.Int64("telegram_id", telegramID),
				zap.Any("meta", richErr.Meta()),
			)
		}
		return
	}

	req := orderparams.GetByUserIdRequest{
		UserID: userID,
	}

	resp, oErr := h.orderService.GetByUserID(ctx, req)
	if oErr != nil {
		logger.Logger.Error("failed to get user orders",
			zap.String("op", op),
			zap.Uint64("user_id", userID),
			zap.Error(err),
		)
		_ = h.messenger.Send(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "❌ خطایی در دریافت تاریخچه رخ داد. لطفاً دوباره تلاش کنید.",
		})
		return
	}

	if len(resp.Orders) == 0 {
		_ = h.messenger.Send(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "📜 شما هنوز هیچ سفارش یا تراکنشی ثبت نکرده‌اید.",
		})
		return
	}

	var sb strings.Builder
	sb.WriteString("📜 <b>آخرین سفارشات شما:</b>\n\n")

	limit := 10
	if len(resp.Orders) < limit {
		limit = len(resp.Orders)
	}

	for i := 0; i < limit; i++ {
		order := resp.Orders[i]
		sb.WriteString(fmt.Sprintf(
			"🔹 <b>سفارش #%d</b>\n"+
				"   📊 وضعیت: %s\n"+
				"   💰 مبلغ: <code>%s</code>\n"+
				"   📅 تاریخ: %s\n\n",
			order.ID,
			order.Status,
			order.Amount.String(),
			order.CreatedAt.Format("2006-01-02 15:04"),
		))
	}

	_ = h.messenger.Send(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      sb.String(),
		ParseMode: models.ParseModeHTML,
	})
}
