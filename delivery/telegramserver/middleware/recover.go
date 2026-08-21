package middleware

import (
	"context"
	"runtime/debug"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"

	"telegram-service-platform/logger"
	"telegram-service-platform/pkg/metrics"
)

func Recover() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			defer func() {
				if r := recover(); r != nil {
					// ۱. ثبت متریک
					metrics.TelegramPanics.Inc()
					metrics.TelegramUpdates.WithLabelValues("panic", "error").Inc()

					// ۲. استخراج اطلاعات update برای لاگ
					updateType, handlerHint, telegramID, chatID := extractUpdateInfo(update)

					// ۳. گرفتن stack trace (فقط تا حد معقول، نه همه‌چیز)
					stack := string(debug.Stack())
					// محدود کردن stack trace برای جلوگیری از نویز در لاگ
					if len(stack) > 4000 {
						stack = stack[:4000] + "\n... (truncated)"
					}

					// ۴. لاگ ساختاریافته با Zap
					logger.Logger.Error("panic recovered in telegram handler",
						zap.String("update_type", updateType),
						zap.String("handler", handlerHint),
						zap.Int64("telegram_id", telegramID),
						zap.Int64("chat_id", chatID),
						zap.Any("panic_value", r),
						zap.String("stack_trace", stack),
					)

					// ۵. ارسال پیام خطای عمومی به کاربر (بدون افشای جزئیات فنی)
					sendPanicMessageToUser(ctx, b, update)
				}
			}()

			next(ctx, b, update)
		}
	}
}

// sendPanicMessageToUser پیام خطای کاربرپسند به کاربر ارسال می‌کند
func sendPanicMessageToUser(ctx context.Context, b *bot.Bot, update *models.Update) {
	errorMsg := "⚠️ خطایی در پردازش درخواست شما رخ داد. لطفاً چند لحظه بعد دوباره تلاش کنید."

	var chatID int64
	if update.Message != nil {
		chatID = update.Message.Chat.ID
	} else if update.CallbackQuery != nil && update.CallbackQuery.Message.Message != nil {
		chatID = update.CallbackQuery.Message.Message.Chat.ID
	} else if update.EditedMessage != nil {
		chatID = update.EditedMessage.Chat.ID
	}

	if chatID == 0 {
		return
	}

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   errorMsg,
	})
}
