package userhandler

import "github.com/go-telegram/bot"

func (h Handler) RegisterRoutes(b *bot.Bot) {
	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/start",
		bot.MatchTypeExact,
		h.start,
	)
}
