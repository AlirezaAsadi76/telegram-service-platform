package priceservice

type Service struct {
	currency    CurrencyProvider
	telegramPrv TelegramProductProvider
	repository  PriceRepository
}

func New(priceRepo PriceRepository, telegramPrv TelegramProductProvider, currency CurrencyProvider) Service {
	return Service{
		currency:    currency,
		telegramPrv: telegramPrv,
		repository:  priceRepo,
	}
}
