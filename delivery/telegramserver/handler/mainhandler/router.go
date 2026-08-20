package mainhandler

import (
	"telegram-service-platform/delivery/telegramserver/callback"
	"telegram-service-platform/delivery/telegramserver/middleware"

	"github.com/go-telegram/bot"
)

func (h *Handler) RegisterRoutes(b *bot.Bot) {

	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/start",
		bot.MatchTypePrefix,
		h.start,
		middleware.Public()...,
	)

	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		callback.ReplyMainMenuText,
		bot.MatchTypeExact,
		h.showMainMenu,
		middleware.Public()...,
	)

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		callback.SMMPlatformPrefix,
		bot.MatchTypePrefix,
		h.selectPlatform,
		middleware.Public()...,
	)

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		callback.SMMCategoryPrefix,
		bot.MatchTypePrefix,
		h.selectCategory,
		middleware.Public()...,
	)

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		callback.SMMServicePrefix,
		bot.MatchTypePrefix,
		h.selectService,
		middleware.Public()...,
	)
}
