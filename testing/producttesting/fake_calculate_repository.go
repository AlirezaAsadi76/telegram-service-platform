package producttesting

import (
	"context"
	"errors"
	"telegram-service-platform/entity/productentity"
	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/service/pricingservice"
	"telegram-service-platform/service/productservice"
)

type calculateSMMPriceRepository struct {
	productservice.Repository

	mapping    *smmentity.SmmMapping
	mappingErr error

	service    *smmentity.SMM
	serviceErr error
}

func (f *calculateSMMPriceRepository) SMMMappingGetByID(
	_ context.Context,
	id int64,
) (*smmentity.SmmMapping, error) {
	if f.mappingErr != nil {
		return nil, f.mappingErr
	}

	if f.mapping == nil ||
		f.mapping.Id != id {
		return nil, errors.New("mapping not found")
	}

	return f.mapping, nil
}

func (f *calculateSMMPriceRepository) SMMServiceGetByD(
	_ context.Context,
	id int64,
) (*smmentity.SMM, error) {
	if f.serviceErr != nil {
		return nil, f.serviceErr
	}

	if f.service == nil ||
		f.service.Id != id {
		return nil, errors.New("service not found")
	}

	return f.service, nil
}

type calculateSMMPriceRepositoryForPricing struct {
	tonUSD   float64
	usdToman float64

	tonErr   error
	tomanErr error
}

func (f *calculateSMMPriceRepositoryForPricing) GetTonUsdPrice(
	_ context.Context,
) (float64, error) {
	return f.tonUSD, f.tonErr
}

func (f *calculateSMMPriceRepositoryForPricing) GetUsdTomanPrice(
	_ context.Context,
) (float64, error) {
	return f.usdToman, f.tomanErr
}

func (f *calculateSMMPriceRepositoryForPricing) GetStarPrice(
	_ context.Context,
) (productentity.StarPrice, error) {
	return productentity.StarPrice{}, nil
}

func (f *calculateSMMPriceRepositoryForPricing) GetPremiumPrices(
	_ context.Context,
) ([]productentity.PremiumPrice, error) {
	return nil, nil
}

func newSMMPricingService(
	tonUSD float64,
	usdToman float64,
) *pricingservice.Service {
	priceRepository := &calculateSMMPriceRepositoryForPricing{
		tonUSD:   tonUSD,
		usdToman: usdToman,
	}

	return pricingservice.New(
		priceRepository,
		pricingservice.Config{
			SMMMinimumUnitPriceToman: "5000",
			SMMMinimumMultiplier:     "1.4",
			SMMMaximumMultiplier:     "4",
			SMMMultiplierScaleToman:  "1466.6666666666666666666666666667",
		},
	)
}
