package priceservice

import (
	"context"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s Service) RefreshCurrency(ctx context.Context) error {
	const Op = "priceservice.RefreshCurrency"
	tonUsd, tErr := s.currency.GetTonUsdPrice(ctx)
	if tErr != nil {
		return richerror.New(Op, tErr).WithKind(richerror.KindDependency).WithMessage(msgerror.ExternalServiceFailed)
	}
	usdIrr, uErr := s.currency.GetUsdTomanPrice(ctx)

	if uErr != nil {
		return richerror.New(Op, uErr).WithKind(richerror.KindDependency).WithMessage(msgerror.ExternalServiceFailed)
	}
	stErr := s.repository.SetTonUsdPrice(ctx, tonUsd, s.config.CurrencyTTL)
	if stErr != nil {
		return richerror.New(Op, stErr).
			WithKind(richerror.KindInfrastructure).
			WithMessage(msgerror.CacheWriteFailed)
	}
	suErr := s.repository.SetUsdTomanPrice(ctx, usdIrr, s.config.CurrencyTTL)
	if suErr != nil {
		return richerror.New(Op, stErr).
			WithKind(richerror.KindInfrastructure).
			WithMessage(msgerror.CacheWriteFailed)

	}

	return nil
}
