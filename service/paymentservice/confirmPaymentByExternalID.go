package paymentservice

import (
	"context"

	"telegram-service-platform/params/paymentparams"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s *Service) ConfirmPaymentByExternalID(ctx context.Context, req paymentparams.ConfirmPaymentByExternalIDRequest) (*paymentparams.ConfirmPaymentResponse, error) {
	const Op = "paymentservice.confirmpaymentbyexternalid"

	payment, err := s.repo.GetByExternalID(ctx, req.ExternalID)
	if err != nil {
		return nil, richerror.New(Op, err).
			WithKind(richerror.KindNotFound).
			WithCode(richerror.CodePaymentNotFound).
			WithMessage(msgerror.PaymentNotFound)
	}

	return s.ConfirmPayment(
		ctx,
		paymentparams.ConfirmPaymentRequest{
			PaymentID:    payment.ID,
			ExternalID:   req.ExternalID,
			CallbackData: req.CallbackData,
		},
	)
}
