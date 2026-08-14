package paymentservice

import (
	"context"
	"telegram-service-platform/params/paymentparams"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s *Service) GetExpired(ctx context.Context) (paymentparams.ExpiredResponse, error) {

	const Op = "paymentservice.getExpired"
	payments, pErr := s.repo.GetPending(ctx)
	if pErr != nil {
		return paymentparams.ExpiredResponse{}, richerror.New(Op, pErr).WithKind(richerror.KindNotFound).WithMessage(msgerror.PaymentNotFound)
	}
	return paymentparams.ExpiredResponse{
		Payments: payments,
	}, nil
}
