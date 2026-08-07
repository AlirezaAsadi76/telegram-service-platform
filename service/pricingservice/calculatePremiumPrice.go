package pricingservice

import (
	"context"
	"telegram-service-platform/entity/productentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s Service) CalculatePremiumPrices(ctx context.Context) (map[uint8]productentity.Price, error) {
	const Op = "pricingservice.CalculatePremiumPrices"

	pricePlans, pErr := s.priceRepository.GetPremiumPrices(ctx)
	if pErr != nil {
		return nil, richerror.New(Op, pErr).WithKind(richerror.KindInfrastructure).WithMessage(msgerror.CacheReadFailed)
	}

	premiumPriceMap := make(map[uint8]productentity.Price, len(pricePlans))

	for _, price := range pricePlans {
		pr, scErr := s.CalculatePrice(ctx, price.PriceUSD)
		if scErr != nil {
			return nil, richerror.New(Op, scErr).
				WithKind(richerror.KindUnexpected).
				WithMessage(msgerror.Unexpected)
		}
		premiumPriceMap[price.Months] = pr
	}
	return premiumPriceMap, nil
}
