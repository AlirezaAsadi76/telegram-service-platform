package dispatcher

import (
	"context"
	"telegram-service-platform/entity"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type ConversationHandler interface {
	CanHandle(ctx context.Context, telegramID entity.TelegramId) bool
	Handle(ctx context.Context, b *bot.Bot, update *models.Update)
}
