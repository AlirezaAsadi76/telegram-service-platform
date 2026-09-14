package paymentservice

import (
	"context"

	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/params/paymentparams"
)

func (s *Service) StartPayment(ctx context.Context, req paymentparams.StartPaymentRequest) (*paymentparams.InitiateResponse, error) {
	intent, err := s.CreateIntent(ctx,
		paymentparams.CreateIntentRequest{
			OrderID:        req.OrderID,
			UserID:         req.UserID,
			Method:         req.Method,
			Amount:         req.Amount,
			Currency:       req.Currency,
			IdempotencyKey: req.IdempotencyKey,
		},
	)
	if err != nil {
		return nil, err
	}

	if intent.Status != paymententity.PaymentStatusCreating {
		payment, gErr := s.repo.GetByID(ctx, intent.PaymentID)
		if gErr != nil {
			return nil, gErr
		}

		return &paymentparams.InitiateResponse{
			PaymentID:  payment.ID,
			Status:     payment.Status,
			ExternalID: payment.ExternalID,
			PaymentURL: payment.PaymentURL,
		}, nil
	}

	return s.Initiate(
		ctx,
		paymentparams.InitiateRequest{
			PaymentID:   intent.PaymentID,
			CallbackURL: req.CallbackURL,
			Description: req.Description,
		},
	)
}
