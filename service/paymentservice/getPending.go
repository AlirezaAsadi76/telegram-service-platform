package paymentservice

import (
	"context"
	"telegram-service-platform/params/paymentparams"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s *Service) GetPending(ctx context.Context) (paymentparams.PendingResponse, error) {

	const Op = "paymentservice.GetPending"
	payments, pErr := s.repo.GetPending(ctx)
	if pErr != nil {
		return paymentparams.PendingResponse{}, richerror.New(Op, pErr).WithKind(richerror.KindNotFound).WithMessage(msgerror.PaymentNotFound)
	}
	return paymentparams.PendingResponse{
		Payments: payments,
	}, nil
}
