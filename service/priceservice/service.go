package priceservice

import (
	"context"
	"time"
)

type Provider interface {
	GetTonPrice(ctx context.Context) (float64, error)
}

type Repository interface {
	SetTonPrice(ctx context.Context, price float64, expTime time.Duration) error
	GetTonPrice(ctx context.Context) (float64, error)
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
