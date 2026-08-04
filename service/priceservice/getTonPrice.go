package priceservice

import (
	"context"
	"telegram-service-platform/params"

	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s Service) GetTonPrice(
	ctx context.Context,
) (params.GetTonPriceResponse, error) {

	const Op = "priceservice.GetTonPrice"

	price, err := s.repository.GetTonPrice(ctx)

	if err == nil && price > 0 {
		return params.GetTonPriceResponse{Price: price}, nil
	}

	price, err = s.provider.GetTonPrice(ctx)

	if err != nil {
		return 0,
			richerror.New(
				Op,
				err,
			).
				WithKind(
					richerror.KindUnexpected,
				).
				WithMessage(
					msgerror.Unexpected,
				)
	}

	if err := s.repository.SetTonPrice(
		ctx,
		price,
	); err != nil {

		// cache failure should not break business flow
	}

	return params.GetTonPriceResponse{Price: price}, nil
}
