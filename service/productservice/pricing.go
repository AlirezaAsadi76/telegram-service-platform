package productservice

import (
	"context"
	"telegram-service-platform/entity"

	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s Service) calculatePrice(
	ctx context.Context,
	usd float64,
) (entity.ProductPrice, error) {

	const Op = "productservice.calculatePrice"

	tonPrice, guErr := s.provider.GetTonUsdPrice(ctx)
	if guErr != nil {
		return entity.ProductPrice{},
			richerror.New(Op, guErr).
				WithKind(richerror.KindUnexpected).
				WithMessage(msgerror.Unexpected)
	}

	tomanPrice, gtErr := s.provider.GetUsdTomanPrice(ctx)
	if gtErr != nil {
		return entity.ProductPrice{},
			richerror.New(Op, gtErr).
				WithKind(richerror.KindUnexpected).
				WithMessage(msgerror.Unexpected)
	}

	return entity.ProductPrice{
		USD:   usd,
		USDT:  usd,
		TON:   usd / tonPrice,
		Toman: usd * tomanPrice,
	}, nil
}
