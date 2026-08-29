package walletservice

import (
	"context"
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/walletentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

// TODO- get must be check in validation
func (s *Service) GetOrCreate(ctx context.Context, userID uint64, currency entity.Currency) (*walletentity.Wallet, error) {
	const Op = "walletservice.GetOrCreate"

	wallet, gErr := s.repo.GetByUserID(ctx, userID)
	if gErr == nil {
		return wallet, nil
	}
	wallet = &walletentity.Wallet{
		UserID:   userID,
		Balance:  entity.Amount{},
		Currency: currency,
		Version:  1,
	}
	if err := s.repo.Create(ctx, wallet); err != nil {
		return nil, richerror.New(Op, err).WithKind(richerror.KindCreateFailed).WithMessage(msgerror.InternalServerError)
	}
	return wallet, nil
}
