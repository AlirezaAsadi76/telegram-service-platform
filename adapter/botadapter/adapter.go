package botadapter

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Adapter struct {
	client *bot.Bot
}

func New(config Config) *Adapter {
	client, err := bot.New(
		config.Token,
		bot.WithDefaultHandler(
			func(ctx context.Context, b *bot.Bot, update *models.Update) {
				log.Printf(
					"unknown update: %+v",
					update,
				)
			},
		),
		//bot.WithDebug(),
	)

	if err != nil {
		panic(err)
	}
	return &Adapter{client: client}
}

func (adapter *Adapter) Client() *bot.Bot {
	return adapter.client
}
