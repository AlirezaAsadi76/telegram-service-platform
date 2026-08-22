package telegramserver

import (
	"telegram-service-platform/adapter/botadapter"
	"telegram-service-platform/delivery/telegramserver/dispatcher"
	"telegram-service-platform/delivery/telegramserver/handler"

	"github.com/go-telegram/bot"
)

type Bot struct {
	client     *bot.Bot
	handlers   []handler.Handler
	dispatcher *dispatcher.MessageDispatcher
}

func New(telegramBot *botadapter.Adapter, dispatcher *dispatcher.MessageDispatcher, handlers ...handler.Handler) (*Bot, error) {

	server := &Bot{
		client:     telegramBot.Client(),
		handlers:   handlers,
		dispatcher: dispatcher,
	}

	server.registerRoutes()

	return server, nil
}
