package priceservice

import (
	"context"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s Service) RefreshPricePremium(ctx context.Context) error {
	const Op = "priceservice.RefreshPricePremium"

	starsPricePrv, gErr := s.telegramPrv.GetPremiumPlans(ctx)
	if gErr != nil {
		return richerror.New(Op, gErr).WithKind(richerror.KindDependency).WithMessage(msgerror.ExternalServiceFailed)
	}
	sErr := s.repository.SetPremiumPrices(ctx, starsPricePrv, s.config.PremiumPriceTTL)
	if sErr != nil {
		return richerror.New(Op, sErr).
			WithKind(richerror.KindInfrastructure).
			WithMessage(msgerror.CacheWriteFailed)
	}
	return nil

}
