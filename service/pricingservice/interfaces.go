package pricingservice

import (
	"context"
)

type PriceRepository interface {
	GetTonUsdPrice(ctx context.Context) (float64, error)
	GetUsdTomanPrice(ctx context.Context) (float64, error)
}
