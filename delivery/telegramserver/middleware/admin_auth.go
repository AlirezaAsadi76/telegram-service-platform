package middleware

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"

	"telegram-service-platform/logger"
)

type AdminAuthConfig struct {
	TelegramIDs []int64
}

func AdminAuth(cfg AdminAuthConfig) Middleware {
	adminMap := make(map[int64]bool)
	for _, id := range cfg.TelegramIDs {
		adminMap[id] = true
	}

	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			var telegramID int64

			if update.Message != nil && update.Message.From != nil {
				telegramID = update.Message.From.ID
			} else if update.CallbackQuery != nil && update.CallbackQuery.From != nil {
				telegramID = update.CallbackQuery.From.ID
			}

			if !adminMap[telegramID] {
				logger.Logger.Warn("unauthorized admin command attempt",
					zap.Int64("telegram_id", telegramID),
				)

				if update.Message != nil {
					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID: update.Message.Chat.ID,
						Text:   "⛔ دسترسی غیرمجاز. شما ادمین نیستید.",
					})
				}
				return
			}

			next(ctx, b, update)
		}
	}
}
