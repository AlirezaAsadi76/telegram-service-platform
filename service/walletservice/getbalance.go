package walletservice

import (
	"context"
	"telegram-service-platform/params/walletparam"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s *Service) GetBalance(ctx context.Context, req walletparam.GetBalanceRequest) (*walletparam.GetBalanceResponse, error) {
	const Op = "walletservice.GetBalance"

	wallet, gfErr := s.repo.GetByUserID(ctx, req.UserID)
	if gfErr != nil {
		return nil, richerror.New(Op, gfErr).WithKind(richerror.KindNotFound).WithMessage(msgerror.WalletNotFound)
	}
	return &walletparam.GetBalanceResponse{
		Balance:  wallet.Balance,
		Currency: wallet.Currency,
	}, nil
}
