package productservice

import "context"

type CurrencyProvider interface {
	GetTonUsdPrice(ctx context.Context) (float64, error)
	GetUsdTomanPrice(ctx context.Context) (float64, error)
}

type ProductProvider interface {
	GetStarsPrice(ctx context.Context) (StarsPrice, error)
	GetPremiumPlans(ctx context.Context) ([]PremiumPrice, error)
}
