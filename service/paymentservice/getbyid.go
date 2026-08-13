package paymentservice

import (
	"context"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s *Service) GetById(ctx context.Context, paymentID uint64) (*paymententity.Payment, error) {
	const Op = "paymentservice.getById"
	payment, iErr := s.repo.GetByID(ctx, paymentID)
	if iErr != nil {
		return nil, richerror.New(Op, iErr).WithKind(richerror.KindNotFound).WithMessage(msgerror.PaymentNotFound)
	}
	return payment, nil
}
