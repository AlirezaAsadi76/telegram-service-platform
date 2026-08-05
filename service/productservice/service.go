package productservice

import (
	"context"
	"telegram-service-platform/entity"
	"time"
)

type Provider interface {
	GetTonPrice(ctx context.Context) (float64, error)
	GetUsdToTomanPrice(ctx context.Context) (float64, error)
}

type Repository interface {
	GetStarsPackage(ctx context.Context) ([]entity.StarPackage, error)
	GetPremiumPlans(ctx context.Context) ([]entity.PremiumPlan, error)
	GetAdsPlans(ctx context.Context) ([]entity.AdsPlan, error)
}

type Config struct {
	ExpTime time.Duration `koanf:"expTime"`
}
type Service struct {
	provider   Provider
	repository Repository
	config     Config
}

func New(config Config, provider Provider, repository Repository) Service {
	return Service{
		provider:   provider,
		repository: repository,
		config:     config,
	}
}
