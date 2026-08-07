package productservice

import (
	"time"
)

type Config struct {
	PriceCacheTTL time.Duration `koanf:"priceCacheTTL"`
}
type Service struct {
	repository Repository
	pricingSVc PricingService
	config     Config
}

func New(config Config, pricingSVc PricingService, repository Repository) Service {
	return Service{

		repository: repository,
		config:     config,
		pricingSVc: pricingSVc,
	}
}
