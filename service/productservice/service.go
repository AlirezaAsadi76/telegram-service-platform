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
	pricingSVc PricingService
	config     Config
}

func New(config Config, pricingSVc PricingService, repository Repository, priceRepo PriceRepository) Service {
	return Service{

		repository: repository,
		priceRepo:  priceRepo,
		config:     config,
		pricingSVc: pricingSVc,
	}
}
