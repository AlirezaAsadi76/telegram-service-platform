package mainhandler

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
	"telegram-service-platform/delivery/telegramserver/keyboard"
	"telegram-service-platform/logger"
	"telegram-service-platform/pkg/metrics"
)

func (h *Handler) start(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID

	// ثبت متریک
	metrics.SMMBotOrdersTotal.WithLabelValues("start_command", "triggered").Inc()

	// خواندن پلتفرم‌ها و دسته‌بندی‌های تلگرام از دیتابیس
	platforms, err := h.productService.GetDistinctPlatforms(ctx)
	if err != nil {
		logger.Logger.Error("failed to get platforms", zap.Error(err))
		h.messenger.SendText(ctx, chatID, " خطا در بارگذاری منو. لطفاً دوباره تلاش کنید.", nil)
		return
	}

	telegramCategories, err := h.productService.GetDistinctCategoriesByPlatform(ctx, "telegram")
	if err != nil {
		logger.Logger.Error("failed to get telegram categories", zap.Error(err))
		h.messenger.SendText(ctx, chatID, " خطا در بارگذاری منو. لطفاً دوباره تلاش کنید.", nil)
		return
	}

	// ساخت پیام خوش‌آمدگویی
	welcomeMsg := fmt.Sprintf(
		"سلام %s! 👋\nبه پلتفرم خدمات SMM خوش آمدید.\nلطفاً از منوی زیر انتخاب کنید:",
		update.Message.From.FirstName,
	)

	// ارسال کیبورد اصلی (Inline)
	inlineKeyboard := keyboard.MainMenu(platforms, telegramCategories)
	h.messenger.SendText(ctx, chatID, welcomeMsg, inlineKeyboard)

	// ارسال کیبورد ریپلای دائمی
	replyKeyboard := keyboard.ReplyMainMenu()
	h.messenger.SendText(ctx, chatID, " دکمه‌های دسترسی سریع پایین صفحه همیشه در دسترس شما هستند:", replyKeyboard)
}

func (h *Handler) showMainMenu(ctx context.Context, b *bot.Bot, update *models.Update) {
	// بازنشانی منو
	h.start(ctx, b, update)
}
