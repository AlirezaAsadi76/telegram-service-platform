package middleware

import (
	"context"
	"telegram-service-platform/entity"
	"telegram-service-platform/logger"
	"telegram-service-platform/params"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

type ActivityTracker interface {
	TrackActivity(ctx context.Context, req params.TrackUserActivityRequest) (params.TrackUserActivityResponse, error)
}

func ActivityTracking(tracker ActivityTracker) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			var telegramID entity.TelegramId

			if update.Message != nil && update.Message.From != nil {
				telegramID = entity.TelegramId(update.Message.From.ID)
			} else if update.CallbackQuery != nil && update.CallbackQuery.Message.Message.From != nil {
				telegramID = entity.TelegramId(update.CallbackQuery.From.ID)
			}

			if telegramID > 0 {
				if _, err := tracker.TrackActivity(ctx, params.TrackUserActivityRequest{TelegramID: telegramID}); err != nil {

					logger.Logger.Warn("failed to track user activity",
						zap.Int64("telegram_id", telegramID.Int64()),
						zap.Error(err),
					)
				}
			}

			next(ctx, b, update)
		}
	}
}
