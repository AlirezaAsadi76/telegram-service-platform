package paymentservice

import (
	"context"
	"telegram-service-platform/pkg/msgerror"

	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/params/paymentparams"
	"telegram-service-platform/pkg/richerror"
)

func (s *Service) Create(ctx context.Context, req paymentparams.CreateRequest) (*paymentparams.CreateResponse, error) {
	const Op = "paymentservice.create"

	provider := s.getProvider(req.Method)

	// Create external payment
	providerReq := paymentparams.CreateRequest{
		Amount:   req.Amount,
		Currency: req.Currency,
	}
	providerResp, pcErr := provider.Create(ctx, providerReq)
	if pcErr != nil {
		return nil, richerror.New(Op, pcErr).WithKind(richerror.KindCreateFailed).WithMessage(msgerror.ProviderCreateFailed)
	}

	// Save to DB
	payment := &paymententity.Payment{
		OrderID:        req.OrderID,
		UserID:         req.UserID,
		Method:         req.Method,
		Amount:         req.Amount,
		Currency:       req.Currency,
		Status:         paymententity.PaymentStatusPending,
		ExternalID:     providerResp.ExternalID,
		IdempotencyKey: req.IdempotencyKey,
	}
	if err := s.repo.Create(ctx, payment); err != nil {
		return nil, richerror.New(Op, pcErr).WithKind(richerror.KindQueryFailure).WithMessage(msgerror.QueryFailed)
	}

	return &paymentparams.CreateResponse{
		PaymentID:  payment.ID,
		PaymentURL: providerResp.PaymentURL,
		ExternalID: providerResp.ExternalID,
	}, nil
}
