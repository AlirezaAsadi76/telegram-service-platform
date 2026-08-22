package telegramserver

import (
	"log"
	"telegram-service-platform/delivery/telegramserver/middleware"

	"github.com/go-telegram/bot"
)

func (b *Bot) registerRoutes() {
	log.Println("register routes:", len(b.handlers))
	for _, h := range b.handlers {
		h.RegisterRoutes(b.client)
	}

	if b.dispatcher != nil {
		b.client.RegisterHandler(
			bot.HandlerTypeMessageText,
			"",
			bot.MatchTypePrefix,
			b.dispatcher.ProcessTextMessage,
			middleware.Public()...,
		)
		log.Println("conversation dispatcher registered as catch-all handler")
	}
}
