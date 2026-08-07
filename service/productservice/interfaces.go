package productservice

import (
	"context"
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/productentity"
)

type Repository interface {
	GetStarPlans(ctx context.Context) ([]entity.StarPackage, error)
	GetPremiumPlans(ctx context.Context) ([]entity.PremiumPlan, error)
	GetAdsPlans(ctx context.Context) ([]entity.AdsPlan, error)
}

type PriceRepository interface {
	GetStarPrice(ctx context.Context) (productentity.StarPrice, error)
	GetPremiumPrices(ctx context.Context) ([]productentity.PremiumPrice, error)
	GetTonUsdPrice(ctx context.Context) (float64, error)
	GetUsdTomanPrice(ctx context.Context) (float64, error)
}
