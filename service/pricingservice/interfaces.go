package pricingservice

import (
	"context"
	"telegram-service-platform/entity/productentity"
)

type PriceRepository interface {
	GetTonUsdPrice(ctx context.Context) (float64, error)
	GetUsdTomanPrice(ctx context.Context) (float64, error)
	GetStarPrice(ctx context.Context) (productentity.StarPrice, error)
	GetPremiumPrices(ctx context.Context) ([]productentity.PremiumPrice, error)
}
