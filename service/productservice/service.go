package productservice

import (
	"time"
)

type Config struct {
	PriceCacheTTL time.Duration `koanf:"priceCacheTTL"`
}
type Service struct {
	repository Repository
	priceRepo  PriceRepository
	config     Config
}

func New(config Config, repository Repository, priceRepo PriceRepository) Service {
	return Service{

		repository: repository,
		priceRepo:  priceRepo,
		config:     config,
	}
}
