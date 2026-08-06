package priceservice

type Service struct {
	currency    CurrencyProvider
	telegramPrv TelegramProductProvider
}

func New(telegramPrv TelegramProductProvider, currency CurrencyProvider) Service {
	return Service{
		currency:    currency,
		telegramPrv: telegramPrv,
	}
}
