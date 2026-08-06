package priceservice

import "context"

type CurrencyProvider interface {
	GetTonUsdPrice(ctx context.Context) (float64, error)
	GetUsdTomanPrice(ctx context.Context) (float64, error)
}

type Service struct {
	currency CurrencyProvider
}

func New(currency CurrencyProvider) Service {
	return Service{
		currency: currency,
	}
}
