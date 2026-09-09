package paymentservice

import (
	"context"
	"time"

	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/params/paymentparams"
	"telegram-service-platform/params/paymentproviderparams"
	"telegram-service-platform/pkg/metrics"
	"telegram-service-platform/pkg/richerror"
)

func (s *Service) Initiate(ctx context.Context, req paymentparams.InitiateRequest) (*paymentparams.InitiateResponse, error) {
	const op = "paymentservice.initiate"

	start := time.Now()

	defer func() {
		metrics.PaymentInitiationDuration.
			Observe(time.Since(start).Seconds())
	}()

	payment, err := s.repo.GetByID(ctx, req.PaymentID)
	if err != nil {
		metrics.PaymentInitiationResult.
			WithLabelValues("unknown", "load_failed").
			Inc()

		return nil, richerror.New(op, err).
			WithKind(richerror.KindQueryFailure).
			WithCode(richerror.CodePaymentLoadFailed)
	}

	if payment.Status != paymententity.PaymentStatusCreating {
		metrics.PaymentInitiationResult.
			WithLabelValues(string(payment.Method), "invalid_state").
			Inc()

		return nil, richerror.New(op, nil).
			WithKind(richerror.KindConflict).
			WithCode(richerror.CodePaymentInvalidState)
	}

	provider := s.getProvider(payment.Method)

	if provider == nil {
		metrics.PaymentInitiationResult.
			WithLabelValues(string(payment.Method), "provider_unavailable").
			Inc()

		return nil, richerror.New(op, nil).
			WithKind(richerror.KindDependencyFailure).
			WithCode(richerror.CodePaymentProviderUnavailable)
	}

	providerReq := paymentproviderparams.CreateRequest{
		PaymentID:   payment.ID,
		OrderID:     payment.OrderID,
		Amount:      payment.Amount,
		Currency:    payment.Currency,
		CallbackURL: req.CallbackURL,
		Description: req.Description,
	}

	providerResp, cErr := provider.Create(ctx, providerReq)
	if cErr != nil {
		if richerror.IsCode(
			cErr,
			richerror.CodePaymentProviderRejected,
		) {
			if err := s.repo.UpdateStatus(
				ctx,
				payment.ID,
				paymententity.PaymentStatusFailed,
			); err != nil {
				metrics.PaymentInitiationResult.
					WithLabelValues(
						string(payment.Method),
						"persist_failed",
					).
					Inc()

				return nil, richerror.New(op, err).
					WithKind(richerror.KindQueryFailure).
					WithCode(
						richerror.CodePaymentInitiationUpdateFailed,
					)
			}

			metrics.PaymentInitiationResult.
				WithLabelValues(
					string(payment.Method),
					"rejected",
				).
				Inc()

			return nil, cErr
		}

		if err := s.repo.UpdateStatus(
			ctx,
			payment.ID,
			paymententity.PaymentStatusUnknown,
		); err != nil {
			metrics.PaymentInitiationResult.
				WithLabelValues(
					string(payment.Method),
					"persist_failed",
				).
				Inc()

			return nil, richerror.New(op, err).
				WithKind(richerror.KindQueryFailure).
				WithCode(
					richerror.CodePaymentInitiationUpdateFailed,
				)
		}

		metrics.PaymentInitiationResult.
			WithLabelValues(
				string(payment.Method),
				"unknown",
			).
			Inc()

		return nil, cErr
	}

	if providerResp.ExternalID == "" {
		if err := s.repo.UpdateStatus(
			ctx,
			payment.ID,
			paymententity.PaymentStatusUnknown,
		); err != nil {
			metrics.PaymentInitiationResult.
				WithLabelValues(
					string(payment.Method),
					"persist_failed",
				).
				Inc()

			return nil, richerror.New(op, err).
				WithKind(richerror.KindQueryFailure).
				WithCode(
					richerror.CodePaymentInitiationUpdateFailed,
				)
		}

		metrics.PaymentInitiationResult.
			WithLabelValues(
				string(payment.Method),
				"invalid_response",
			).
			Inc()

		return nil, richerror.New(op, nil).
			WithKind(richerror.KindExternalAPI).
			WithCode(
				richerror.CodePaymentProviderInvalidResponse,
			)
	}

	if err := s.repo.MarkInitiated(
		ctx,
		payment.ID,
		paymententity.PaymentStatusPending,
		providerResp.ExternalID,
		providerResp.PaymentURL,
	); err != nil {
		metrics.PaymentInitiationResult.
			WithLabelValues(string(payment.Method), "persist_failed").
			Inc()

		// IMPORTANT:
		// Provider may already have created the payment.
		// Do NOT mark it FAILED.
		return nil, richerror.New(op, err).
			WithKind(richerror.KindQueryFailure).
			WithCode(richerror.CodePaymentInitiationUpdateFailed)
	}

	metrics.PaymentInitiationResult.
		WithLabelValues(string(payment.Method), "created").
		Inc()

	return &paymentparams.InitiateResponse{
		PaymentID:  payment.ID,
		Status:     paymententity.PaymentStatusPending,
		PaymentURL: providerResp.PaymentURL,
		ExternalID: providerResp.ExternalID,
	}, nil
}
