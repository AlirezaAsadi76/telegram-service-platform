package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"

	"telegram-service-platform/logger"
	"telegram-service-platform/pkg/metrics"
)

func Logger() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			start := time.Now()

			// استخراج اطلاعات update
			updateType, handlerHint, telegramID, chatID := extractUpdateInfo(update)

			// ثبت متریک ورودی
			metrics.TelegramUpdates.WithLabelValues(updateType, "received").Inc()

			// اجرای handler
			next(ctx, b, update)

			// latency دقیق (بعد از اجرای handler)
			duration := time.Since(start)

			// ثبت متریک duration
			metrics.TelegramHandlerDuration.WithLabelValues(updateType, handlerHint).Observe(duration.Seconds())

			// لاگ ساختاریافته با Zap
			logger.Logger.Info("telegram update processed",
				zap.String("update_type", updateType),
				zap.String("handler", handlerHint),
				zap.Int64("telegram_id", telegramID),
				zap.Int64("chat_id", chatID),
				zap.Duration("duration", duration),
			)
		}
	}
}

// extractUpdateInfo اطلاعات کلیدی را از update استخراج می‌کند
func extractUpdateInfo(update *models.Update) (updateType, handlerHint string, telegramID, chatID int64) {
	if update == nil {
		return "unknown", "unknown", 0, 0
	}

	// Callback Query (دکمه‌های Inline)
	if update.CallbackQuery != nil {
		telegramID = update.CallbackQuery.From.ID
		if update.CallbackQuery.Message.Message != nil {
			chatID = update.CallbackQuery.Message.Message.Chat.ID
		}
		data := update.CallbackQuery.Data
		handlerHint = extractHandlerHintFromCallback(data)
		return "callback", handlerHint, telegramID, chatID
	}

	// Message (پیام متنی یا دستور)
	if update.Message != nil {
		if update.Message.From != nil {
			telegramID = update.Message.From.ID
		}
		chatID = update.Message.Chat.ID
		text := update.Message.Text

		if strings.HasPrefix(text, "/") {
			// دستور مثل /start, /recharge
			parts := strings.Fields(text)
			if len(parts) > 0 {
				handlerHint = parts[0] // مثلاً "/start"
			}
			return "command", handlerHint, telegramID, chatID
		}
		return "message", "text", telegramID, chatID
	}

	// Edited Message
	if update.EditedMessage != nil {
		if update.EditedMessage.From != nil {
			telegramID = update.EditedMessage.From.ID
		}
		chatID = update.EditedMessage.Chat.ID
		return "edited_message", "text", telegramID, chatID
	}

	// Inline Query
	if update.InlineQuery != nil {
		telegramID = update.InlineQuery.From.ID
		return "inline_query", update.InlineQuery.Query, telegramID, 0
	}

	// Chosen Inline Result
	if update.ChosenInlineResult != nil {
		telegramID = update.ChosenInlineResult.From.ID
		return "chosen_inline_result", update.ChosenInlineResult.ResultID, telegramID, 0
	}

	return "unknown", "unknown", 0, 0
}

// extractHandlerHintFromCallback از callback data، hint مربوط به handler را استخراج می‌کند
// مثال: "smm:platform:telegram" → "smm:platform"
// مثال: "product:stars:buy:1" → "product:stars:buy"
func extractHandlerHintFromCallback(data string) string {
	if data == "" {
		return "empty"
	}
	parts := strings.Split(data, ":")
	// دو بخش اول کافی است تا handler را مشخص کنیم
	if len(parts) >= 2 {
		return parts[0] + ":" + parts[1]
	}
	return parts[0]
}
