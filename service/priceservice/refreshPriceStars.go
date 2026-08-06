package priceservice

import (
	"context"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s Service) RefreshPriceStars(ctx context.Context) error {
	const Op = "priceservice.refreshPriceStars"
	starsPricePrv, gErr := s.telegramPrv.GetStarPrice(ctx)
	if gErr != nil {
		return richerror.New(Op, gErr).WithKind(richerror.KindDependency).WithMessage(msgerror.ExternalServiceFailed)
	}
	sErr := s.repository.SetStarPrice(ctx, starsPricePrv, s.config.StarsPriceTTL)
	if sErr != nil {
		return richerror.New(Op, sErr).
			WithKind(richerror.KindInfrastructure).
			WithMessage(msgerror.CacheWriteFailed)
	}
	return nil

}
