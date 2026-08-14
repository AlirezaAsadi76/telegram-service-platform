package paymentservice

import (
	"context"
	"telegram-service-platform/params/paymentparams"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s *Service) UpdateStatus(ctx context.Context, request paymentparams.UpdateStatusRequest) error {

	const Op = "paymentservice.UpdateStatus"
	pErr := s.repo.UpdateStatus(ctx, request.PymentId, request.PaymentStatus)
	if pErr != nil {
		return richerror.New(Op, pErr).WithKind(richerror.KindNotFound).WithMessage(msgerror.PaymentNotFound)
	}
	return nil
}
