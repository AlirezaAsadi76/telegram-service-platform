package pricingservice

type Service struct {
	priceRepository PriceRepository
	config          Config
}

func New(priceRepo PriceRepository, cfg Config) *Service {
	return &Service{
		priceRepository: priceRepo,
		config:          cfg,
	}
}
