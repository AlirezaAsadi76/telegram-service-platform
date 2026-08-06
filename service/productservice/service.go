package productservice

import (
	"context"
	"telegram-service-platform/entity"
	"time"
)

type PriceProvider interface {
	GetTonUsdPrice(ctx context.Context) (float64, error)
	GetUsdTomanPrice(ctx context.Context) (float64, error)
}

type Repository interface {
	GetStarPlans(ctx context.Context) ([]entity.StarPackage, error)
	GetPremiumPlans(ctx context.Context) ([]entity.PremiumPlan, error)
	GetAdsPlans(ctx context.Context) ([]entity.AdsPlan, error)
}

type Config struct {
	PriceCacheTTL time.Duration `koanf:"priceCacheTTL"`
}
type Service struct {
	provider   PriceProvider
	repository Repository
	config     Config
}

func New(config Config, provider PriceProvider, repository Repository) Service {
	return Service{
		provider:   provider,
		repository: repository,
		config:     config,
	}
}
