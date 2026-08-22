package productservice

import (
	"telegram-service-platform/repository/redis/rediscatalog"
	"telegram-service-platform/repository/redis/redissmm"
	"telegram-service-platform/service/pricingservice"
	"time"
)

type Config struct {
	PriceCacheTTL time.Duration `koanf:"priceCacheTTL"`
}
type Service struct {
	repository   Repository
	pricingSVc   *pricingservice.Service
	adapter      SMMAdapterInterface
	catalogCache *rediscatalog.CatalogCache
	smmCache     *redissmm.SMMCache
	config       Config
}

func New(config Config, pricingSVc *pricingservice.Service, repository Repository, catalogCache *rediscatalog.CatalogCache, smmCache *redissmm.SMMCache, adapter SMMAdapterInterface) *Service {
	return &Service{

		repository:   repository,
		config:       config,
		pricingSVc:   pricingSVc,
		catalogCache: catalogCache,
		adapter:      adapter,
		smmCache:     smmCache,
	}
}
