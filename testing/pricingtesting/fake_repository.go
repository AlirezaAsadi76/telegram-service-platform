package pricingtesting

import (
	"context"

	"telegram-service-platform/entity/productentity"
	"telegram-service-platform/service/pricingservice"
)

var _ pricingservice.PriceRepository = (*fakePriceRepository)(nil)

type fakePriceRepository struct {
	tonUsdPrice   float64
	usdTomanPrice float64

	tonUsdErr   error
	usdTomanErr error

	tonUsdCalls   int
	usdTomanCalls int
}

func newSMMPricingServiceWithConfig(
	config pricingservice.Config,
	tonUSD float64,
	usdToman float64,
) *pricingservice.Service {
	repository := &fakePriceRepository{
		tonUsdPrice:   tonUSD,
		usdTomanPrice: usdToman,
	}

	return pricingservice.New(
		repository,
		config,
	)
}

func (f *fakePriceRepository) GetTonUsdPrice(
	_ context.Context,
) (float64, error) {
	f.tonUsdCalls++

	return f.tonUsdPrice, f.tonUsdErr
}

func (f *fakePriceRepository) GetUsdTomanPrice(
	_ context.Context,
) (float64, error) {
	f.usdTomanCalls++

	return f.usdTomanPrice, f.usdTomanErr
}

func (f *fakePriceRepository) GetStarPrice(
	_ context.Context,
) (productentity.StarPrice, error) {
	return productentity.StarPrice{}, nil
}

func (f *fakePriceRepository) GetPremiumPrices(
	_ context.Context,
) ([]productentity.PremiumPrice, error) {
	return nil, nil
}
