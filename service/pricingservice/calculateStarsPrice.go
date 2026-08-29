package pricingservice

import (
	"context"
	"telegram-service-platform/entity/productentity"

	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s Service) CalculateStarsPrice(ctx context.Context, amount int64) (productentity.Price, error) {
	const Op = "pricingservice.CalculateStarsPrice"

	priceStars, psErr := s.priceRepository.GetStarPrice(ctx)
	if psErr != nil {
		return productentity.Price{}, richerror.New(Op, psErr).WithKind(richerror.KindInfrastructure).WithMessage(msgerror.CacheReadFailed)
	}

	price, cErr := s.CalculatePrice(ctx, priceStars.PricePerStar.MulInt(amount))
	if cErr != nil {
		return productentity.Price{},
			richerror.New(Op, cErr).
				WithKind(richerror.KindUnexpected).
				WithMessage(msgerror.Unexpected)
	}
	return price, nil
}
