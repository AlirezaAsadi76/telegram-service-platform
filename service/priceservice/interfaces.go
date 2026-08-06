package priceservice

import (
	"context"
	"telegram-service-platform/entity/productentity"
	"time"
)

type CurrencyProvider interface {
	GetTonUsdPrice(ctx context.Context) (float64, error)
	GetUsdTomanPrice(ctx context.Context) (float64, error)
}

type TelegramProductProvider interface {
	GetStarPrice(ctx context.Context) (productentity.StarPrice, error)
	GetPremiumPlans(ctx context.Context) ([]productentity.PremiumPrice, error)
}

type PriceRepository interface {
	SetStarPrice(ctx context.Context, price productentity.StarPrice, expiration time.Duration) error

	SetPremiumPrices(ctx context.Context, prices []productentity.PremiumPrice, expiration time.Duration) error

	SetTonUsdPrice(ctx context.Context, price float64, expiration time.Duration) error

	SetUsdTomanPrice(ctx context.Context, price float64, expiration time.Duration) error
}
