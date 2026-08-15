package pricingservice

type Service struct {
	priceRepository PriceRepository
}

func New(priceRepo PriceRepository) *Service {
	return &Service{
		priceRepository: priceRepo,
	}
}
