package telegramserver

import (
	"telegram-service-platform/adapter/botadapter"
	"telegram-service-platform/delivery/telegramserver/handler"

	"github.com/go-telegram/bot"
)

type Bot struct {
	client   *bot.Bot
	handlers []handler.Handler
}

func New(telegramBot *botadapter.Adapter, handlers ...handler.Handler) (*Bot, error) {

	server := &Bot{
		client:   telegramBot.Client(),
		handlers: handlers,
	}

	server.registerRoutes()

	return server, nil
}
