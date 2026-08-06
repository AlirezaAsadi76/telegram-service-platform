package priceservice

type Service struct {
	currency    CurrencyProvider
	telegramPrv TelegramProductProvider
	repository  PriceRepository
	config      Config
}

func New(cfg Config, priceRepo PriceRepository, telegramPrv TelegramProductProvider, currency CurrencyProvider) Service {
	return Service{
		currency:    currency,
		telegramPrv: telegramPrv,
		repository:  priceRepo,
		config:      cfg,
	}
}
