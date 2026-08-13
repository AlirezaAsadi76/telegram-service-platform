package paymentservice

import (
	"context"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/params/paymentparams"
	"telegram-service-platform/params/paymentproviderparams"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s *Service) Verify(ctx context.Context, req paymentparams.VerifyRequest) (*paymentparams.VerifyResponse, error) {
	const Op = "paymentservice.verify"
	payment, iErr := s.repo.GetByID(ctx, req.PaymentID)
	if iErr != nil {
		return nil, richerror.New(Op, iErr).WithKind(richerror.KindNotFound).WithMessage(msgerror.ProductNotFound)
	}

	// Idempotency: don't verify already processed payments
	if payment.Status == paymententity.PaymentStatusSuccess || payment.Status == paymententity.PaymentStatusFailed {
		return &paymentparams.VerifyResponse{Status: payment.Status}, nil
	}

	provider := s.getProvider(payment.Method)
	providerReq := paymentproviderparams.VerifyRequest{
		ExternalID:   req.ExternalID,
		CallbackData: req.CallbackData,
	}
	providerResp, pvErr := provider.Verify(ctx, providerReq)
	if pvErr != nil {
		return nil, richerror.New(Op, pvErr).WithKind(richerror.KindExternalAPI).WithMessage(msgerror.PaymentVerifyFailed)
	}

	var newStatus paymententity.PaymentStatus
	if providerResp.Status == paymententity.PaymentStatusSuccess {
		newStatus = paymententity.PaymentStatusSuccess
	} else {
		newStatus = paymententity.PaymentStatusFailed
	}

	if err := s.repo.UpdateStatus(ctx, payment.ID, newStatus); err != nil {
		return nil, richerror.New(Op, err).WithKind(richerror.KindInternal).WithMessage(msgerror.InternalServerError)
	}

	return &paymentparams.VerifyResponse{Status: newStatus}, nil
}
