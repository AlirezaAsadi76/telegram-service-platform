package userhandler

import (
	"log"
	"telegram-service-platform/delivery/telegramserver/middleware"

	"github.com/go-telegram/bot"
)

func (h Handler) RegisterRoutes(b *bot.Bot) {
	log.Println("userHandler handler")

	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/start",
		bot.MatchTypeExact,
		h.start,
		middleware.Logger(),
		middleware.Recover(),
	)

}
