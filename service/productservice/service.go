package productservice

import (
	"time"
)

type Config struct {
	PriceCacheTTL time.Duration `koanf:"priceCacheTTL"`
}
type Service struct {
	currencyProvider CurrencyProvider
	telegramProvider TelegramProductProvider
	repository       Repository
	config           Config
}

func New(config Config, telegramProvider TelegramProductProvider, currencyProvider CurrencyProvider, repository Repository) Service {
	return Service{
		currencyProvider: currencyProvider,
		telegramProvider: telegramProvider,
		repository:       repository,
		config:           config,
	}
}
