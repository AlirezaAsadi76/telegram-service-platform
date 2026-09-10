package paymentservice

import (
	"context"
	"fmt"

	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/params/paymentparams"
	"telegram-service-platform/params/paymentproviderparams"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s *Service) ConfirmPayment(ctx context.Context, req paymentparams.ConfirmPaymentRequest) (*paymentparams.ConfirmPaymentResponse, error) {
	const Op = "paymentservice.confirmpayment"

	payment, err := s.repo.GetByID(ctx, req.PaymentID)
	if err != nil {
		return nil, richerror.New(Op, err).
			WithKind(richerror.KindNotFound).
			WithCode(richerror.CodePaymentNotFound).
			WithMessage(msgerror.PaymentNotFound)
	}

	if payment.Status == paymententity.PaymentStatusSuccess {
		return &paymentparams.ConfirmPaymentResponse{
			PaymentID: payment.ID,
			OrderID:   payment.OrderID,
			Status:    payment.Status,
		}, nil
	}

	if payment.Status != paymententity.PaymentStatusPending {
		return nil, richerror.New(Op, nil).
			WithKind(richerror.KindConflict).
			WithCode(richerror.CodePaymentInvalidState).
			WithMessage(msgerror.PaymentConfirmationConflict)
	}

	provider := s.getProvider(payment.Method)
	if provider == nil {
		return nil, richerror.New(
			Op,
			fmt.Errorf("unknown payment method: %s", payment.Method),
		).
			WithKind(richerror.KindInternal).
			WithMessage(msgerror.InternalServerError)
	}

	externalID := req.ExternalID
	if externalID == "" {
		externalID = payment.ExternalID
	}

	providerReq := paymentproviderparams.VerifyRequest{
		PaymentID:    payment.ID,
		ExternalID:   externalID,
		CallbackData: req.CallbackData,
	}

	providerResp, pErr := provider.Verify(ctx, providerReq)
	if pErr != nil {
		return nil, richerror.New(Op, pErr).
			WithKind(richerror.KindExternalAPI).
			WithCode(richerror.CodePaymentVerificationFailed).
			WithMessage(msgerror.PaymentVerifyFailed)
	}

	switch providerResp.Status {
	case paymententity.PaymentStatusSuccess:
		if err := s.paymentConfirmationRepo.Confirm(
			ctx,
			payment.ID,
		); err != nil {
			if richerror.IsKind(err, richerror.KindConflict) {
				if richerror.IsCode(err, richerror.CodePaymentAlreadyConfirmed) {
					return &paymentparams.ConfirmPaymentResponse{
						PaymentID: payment.ID,
						OrderID:   payment.OrderID,
						Status:    payment.Status,
					}, nil
				}
			}
			return nil, richerror.New(Op, err).
				WithKind(richerror.KindInternal).
				WithMessage(msgerror.InternalServerError)
		}

		return &paymentparams.ConfirmPaymentResponse{
			PaymentID: payment.ID,
			OrderID:   payment.OrderID,
			Status:    paymententity.PaymentStatusSuccess,
		}, nil

	case paymententity.PaymentStatusFailed:
		if err := s.paymentConfirmationRepo.Fail(
			ctx,
			payment.ID,
		); err != nil {
			return nil, richerror.New(Op, err).
				WithKind(richerror.KindInternal).
				WithMessage(msgerror.InternalServerError)
		}

		return &paymentparams.ConfirmPaymentResponse{
			PaymentID: payment.ID,
			OrderID:   payment.OrderID,
			Status:    paymententity.PaymentStatusFailed,
		}, nil

	default:
		return nil, richerror.New(Op, nil).
			WithKind(richerror.KindInternal).
			WithCode(richerror.CodePaymentProviderInvalidResponse).
			WithMessage(msgerror.PaymentVerifyFailed)
	}
}
