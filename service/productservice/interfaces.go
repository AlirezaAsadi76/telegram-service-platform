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

type PricingService interface {
	CalculateStarsPrice(ctx context.Context, amount float64) (productentity.Price, error)
	CalculatePremiumPrices(ctx context.Context) (map[uint8]productentity.Price, error)
}
