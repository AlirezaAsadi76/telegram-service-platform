package callbackhandler

import (
	"log"

	"github.com/go-telegram/bot"
)

func (h Handler) RegisterRoutes(b *bot.Bot) {
	log.Println("callback handler")

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		Stars,
		bot.MatchTypePrefix,
		h.callback,
	)

}
