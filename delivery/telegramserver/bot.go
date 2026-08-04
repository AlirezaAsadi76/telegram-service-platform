package telegramserver

import (
	"context"
	"log"
	"telegram-service-platform/delivery/telegramserver/handler/userhandler"
	"telegram-service-platform/service/userservice"
	"telegram-service-platform/validator/uservalidator"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Config struct {
	Token string
	Debug bool
}

type Bot struct {
	client     *bot.Bot
	userRouter userhandler.Router
	config     Config
}

func New(
	cfg Config,
	userRouter userhandler.Router,
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

		client: client,

		userRouter: userRouter,

		config: cfg,
	}

	server.registerRoutes()

	return server, nil
}
