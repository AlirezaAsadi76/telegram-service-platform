package pricingservice

import (
	"context"
	"telegram-service-platform/entity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"

	"github.com/shopspring/decimal"
)

func (s Service) GetTonUsdPrice(ctx context.Context) (entity.Amount, error) {
	const Op = "priceService.GetTonUsdPrice"
	price, err := s.priceRepository.GetTonUsdPrice(ctx)
	if err != nil {
		return entity.Amount{}, richerror.New(Op, err).WithKind(richerror.KindInfrastructure).WithMessage(msgerror.CacheReadFailed)
	}

	return entity.Amount(decimal.NewFromFloat(price)), nil
}
