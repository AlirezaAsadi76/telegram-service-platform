package userhandler

import "github.com/go-telegram/bot"

func (r Router) Register(
	b *bot.Bot,
) {

	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/start",
		bot.MatchTypeExact,
		r.handler.Start,
	)

}
