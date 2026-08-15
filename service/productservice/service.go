package productservice

import (
	"telegram-service-platform/service/pricingservice"
	"time"
)

type Config struct {
	PriceCacheTTL time.Duration `koanf:"priceCacheTTL"`
}
type Service struct {
	repository Repository
	pricingSVc *pricingservice.Service
	config     Config
}

func New(config Config, pricingSVc *pricingservice.Service, repository Repository) *Service {
	return &Service{

		repository: repository,
		config:     config,
		pricingSVc: pricingSVc,
	}
}
