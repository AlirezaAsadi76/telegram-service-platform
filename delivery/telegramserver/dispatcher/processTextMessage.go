package dispatcher

import (
	"context"
	"strings"
	"telegram-service-platform/entity"
	"telegram-service-platform/logger"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

// ProcessTextMessage نقطه ورودی تمام پیام‌های متنی
func (d *MessageDispatcher) ProcessTextMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.Text == "" {
		return
	}

	telegramID := update.Message.From.ID
	text := strings.TrimSpace(update.Message.Text)

	if strings.HasPrefix(text, "/") {
		return
	}

	for _, handler := range d.handlers {
		if handler.CanHandle(ctx, entity.TelegramId(telegramID)) {
			handler.Handle(ctx, b, update)
			return
		}
	}

	logger.Logger.Debug("text message ignored by dispatcher (no active conversation)",
		zap.Int64("telegram_id", telegramID),
		zap.String("text", text),
	)
}
