package mainhandler

import (
	"telegram-service-platform/delivery/telegramserver/callback"

	"github.com/go-telegram/bot"
)

func (h *Handler) RegisterRoutes(b *bot.Bot) {

	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/start",
		bot.MatchTypePrefix,
		h.start,
		h.middlewares()...,
	)

	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		callback.ReplyMainMenuText,
		bot.MatchTypeExact,
		h.showMainMenu,
		h.middlewares()...,
	)

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		callback.SMMPlatformPrefix,
		bot.MatchTypePrefix,
		h.selectPlatform,
		h.middlewares()...,
	)

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		callback.SMMCategoryPrefix,
		bot.MatchTypePrefix,
		h.selectCategory,
		h.middlewares()...,
	)

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		callback.SMMServicePrefix,
		bot.MatchTypePrefix,
		h.selectService,
		h.middlewares()...,
	)
}
