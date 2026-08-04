package priceservice

import "context"

type Provider interface {
	GetTonPrice(ctx context.Context) (float64, error)
}

type Repository interface {
	SetTonPrice(ctx context.Context, price float64) error
	GetTonPrice(ctx context.Context) (float64, error)
}

type Service struct {
	provider   Provider
	repository Repository
}

func New(provider Provider, repository Repository) Service {
	return Service{
		provider:   provider,
		repository: repository,
	}
}
