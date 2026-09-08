package paymentservice

import (
	"context"
	"time"

	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/params/paymentparams"
	"telegram-service-platform/pkg/metrics"
	"telegram-service-platform/pkg/richerror"
)

func (s *Service) CreateIntent(ctx context.Context, req paymentparams.CreateIntentRequest) (*paymentparams.CreateIntentResponse, error) {
	const op = "paymentservice.create_intent"

	start := time.Now()

	defer func() {
		metrics.PaymentIntentDuration.Observe(
			time.Since(start).Seconds(),
		)
	}()

	existingPayment, err := s.repo.GetByIdempotencyKey(ctx, req.IdempotencyKey)
	if err == nil {
		metrics.PaymentIntentResult.WithLabelValues(
			string(existingPayment.Method),
			"existing",
		).Inc()

		return &paymentparams.CreateIntentResponse{
			PaymentID: existingPayment.ID,
			Status:    existingPayment.Status,
		}, nil
	}

	if !richerror.IsKind(err, richerror.KindNotFound) {
		metrics.PaymentIntentResult.WithLabelValues(
			string(req.Method),
			"load_failed",
		).Inc()

		return nil, richerror.New(op, err).
			WithKind(richerror.KindQueryFailure).
			WithCode(richerror.CodePaymentLoadFailed)
	}

	payment := &paymententity.Payment{
		OrderID:        req.OrderID,
		UserID:         req.UserID,
		Method:         req.Method,
		Amount:         req.Amount,
		Currency:       req.Currency,
		Status:         paymententity.PaymentStatusCreating,
		IdempotencyKey: req.IdempotencyKey,
	}

	if err := s.repo.Create(ctx, payment); err != nil {
		
		if richerror.IsKind(err, richerror.KindConflict) {

			existing, getErr := s.repo.GetByIdempotencyKey(
				ctx,
				req.IdempotencyKey,
			)

			if getErr != nil {
				return nil, richerror.New(op, getErr).
					WithKind(richerror.KindQueryFailure).
					WithCode(richerror.CodePaymentLoadFailed)
			}

			metrics.PaymentIntentResult.
				WithLabelValues(string(existing.Method), "conflict_recovered").
				Inc()

			return &paymentparams.CreateIntentResponse{
				PaymentID: existing.ID,
				Status:    existing.Status,
			}, nil
		}

		metrics.PaymentIntentResult.
			WithLabelValues(string(req.Method), "failed").
			Inc()

		return nil, richerror.New(op, err).
			WithKind(richerror.KindCreateFailed).
			WithCode(richerror.CodePaymentIntentCreationFailed)
	}

	metrics.PaymentIntentResult.WithLabelValues(
		string(payment.Method),
		"created",
	).Inc()

	return &paymentparams.CreateIntentResponse{
		PaymentID: payment.ID,
		Status:    payment.Status,
	}, nil
}
