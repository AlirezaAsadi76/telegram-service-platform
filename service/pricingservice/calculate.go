package pricingservice

import (
	"context"
	"errors"
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/productentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"

	"github.com/shopspring/decimal"
)

func (s Service) CalculatePrice(
	ctx context.Context,
	usd entity.Amount,
) (productentity.Price, error) {

	const Op = "pricingservice.CalculatePrice"

	tonPrice, guErr := s.priceRepository.GetTonUsdPrice(ctx)
	if guErr != nil {
		return productentity.Price{},
			richerror.New(Op, guErr).
				WithKind(richerror.KindInfrastructure).
				WithMessage(msgerror.CacheReadFailed)
	}

	tomanPrice, gtErr := s.priceRepository.GetUsdTomanPrice(ctx)
	if gtErr != nil {
		return productentity.Price{},
			richerror.New(Op, gtErr).
				WithKind(richerror.KindInfrastructure).
				WithMessage(msgerror.CacheReadFailed)
	}

	if tonPrice <= 0 {
		return productentity.Price{},
			richerror.New(
				Op,
				errors.New("invalid ton price"),
			).
				WithKind(richerror.KindInvalid).
				WithMessage(msgerror.InvalidPrice)
	}

	tonPriceDecimal := decimal.NewFromFloat(tomanPrice)
	tomanPriceDecimal := decimal.NewFromFloat(tomanPrice)

	return productentity.Price{
		USD:   usd,
		USDT:  usd,
		TON:   usd.Div(entity.Amount(tonPriceDecimal)),
		Toman: usd.Mul(entity.Amount(tomanPriceDecimal)),
	}, nil
}
