package messenger

import (
	"telegram-service-platform/adapter/botadapter"

	"github.com/go-telegram/bot"
)

type Service struct {
	telegramBot *bot.Bot
}

func New(telegramBot *botadapter.Adapter) *Service {

	return &Service{
		telegramBot: telegramBot.Client(),
	}

}
