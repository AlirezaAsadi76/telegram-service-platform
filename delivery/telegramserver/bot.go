package telegramserver

import (
	"context"
	"log"

	"telegram-service-platform/delivery/telegramserver/handler/userhandler"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Config struct {
	Token string `koanf:"token"`
	Debug bool   `koanf:"debug"`
}

type Bot struct {
	client      *bot.Bot
	userHandler userhandler.Handler
	config      Config
}

func New(
	cfg Config,
	userHandler userhandler.Handler,
) (*Bot, error) {

	client, err := bot.New(

		cfg.Token,
		bot.WithDefaultHandler(
			func(
				ctx context.Context,
				b *bot.Bot,
				update *models.Update,
			) {
				log.Printf(
					"unknown update: %+v",
					update,
				)
			},
		),
	)

	if err != nil {
		return nil, err
	}

	server := &Bot{
		client:      client,
		userHandler: userHandler,
		config:      cfg,
	}

	server.registerRoutes()

	return server, nil
}
