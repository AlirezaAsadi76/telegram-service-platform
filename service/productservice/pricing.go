package productservice

import (
	"context"
	"telegram-service-platform/entity/productentity"

	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s Service) calculatePrice(
	ctx context.Context,
	usd float64,
) (productentity.Price, error) {

	const Op = "productservice.calculatePrice"

	tonPrice, guErr := s.priceRepo.GetTonUsdPrice(ctx)
	if guErr != nil {
		return productentity.Price{},
			richerror.New(Op, guErr).
				WithKind(richerror.KindUnexpected).
				WithMessage(msgerror.Unexpected)
	}

	tomanPrice, gtErr := s.priceRepo.GetUsdTomanPrice(ctx)
	if gtErr != nil {
		return productentity.Price{},
			richerror.New(Op, gtErr).
				WithKind(richerror.KindUnexpected).
				WithMessage(msgerror.Unexpected)
	}

	return productentity.Price{
		USD:   usd,
		USDT:  usd,
		TON:   usd / tonPrice,
		Toman: usd * tomanPrice,
	}, nil
}
