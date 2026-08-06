package productservice

import (
	"context"
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/productentity"
)

type CurrencyProvider interface {
	GetTonUsdPrice(ctx context.Context) (float64, error)
	GetUsdTomanPrice(ctx context.Context) (float64, error)
}

type TelegramProductProvider interface {
	GetStarsPrice(ctx context.Context) (productentity.StarPrice, error)
	GetPremiumPlans(ctx context.Context) ([]productentity.PremiumPrice, error)
}

type Repository interface {
	GetStarPlans(ctx context.Context) ([]entity.StarPackage, error)
	GetPremiumPlans(ctx context.Context) ([]entity.PremiumPlan, error)
	GetAdsPlans(ctx context.Context) ([]entity.AdsPlan, error)
}
