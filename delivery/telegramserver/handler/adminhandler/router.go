package adminhandler

import (
	"github.com/go-telegram/bot"
)

func (h *Handler) RegisterRoutes(b *bot.Bot) {

	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/recharge",
		bot.MatchTypePrefix,
		h.Recharge,
		h.middlewares()...,
	)
}
