package smmproviderservice

import (
	"context"
	"time"

	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/entity/providerentity"
	"telegram-service-platform/pkg/metrics"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"

	"telegram-service-platform/logger"

	"go.uber.org/zap"
)

func (s *Service) GetOrderStatus(ctx context.Context, providerID uint64, externalOrderID string) (orderentity.OrderStatus, error) {
	const Op = "smmproviderservice.GetOrderStatus"

	provider, err := s.repo.GetByID(ctx, providerID)
	if err != nil {
		return "", richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryFailed)
	}

	if provider.Type != providerentity.ProviderTypeSMM {
		return "", richerror.New(Op, nil).
			WithKind(richerror.KindDependencyFailure).
			WithCode(richerror.CodeSMMProviderUnavailable).
			WithMessage(msgerror.SMMProviderUnavailable)
	}

	adapter, ok := s.providers[provider.Name]
	if !ok {
		return "", richerror.New(Op, nil).
			WithKind(richerror.KindDependencyFailure).
			WithCode(richerror.CodeSMMProviderUnavailable).
			WithMessage(msgerror.NoAvailableAdapter)
	}

	breaker, okk := s.breakers[provider.Name]
	if !okk {
		return "", richerror.New(Op, nil).
			WithKind(richerror.KindDependencyFailure).
			WithCode(richerror.CodeSMMProviderUnavailable).
			WithMessage(msgerror.SMMProviderUnavailable)
	}

	if !breaker.Allow() {
		return "", richerror.New(Op, nil).
			WithKind(richerror.KindDependencyFailure).
			WithCode(richerror.CodeSMMProviderUnavailable).
			WithMessage(msgerror.SMMProviderUnavailable)
	}

	start := time.Now()

	status, gErr := adapter.GetOrderStatus(ctx, externalOrderID)

	metrics.SMMProviderLatency.
		WithLabelValues(provider.Name).
		Observe(time.Since(start).Seconds())

	if gErr != nil {
		breaker.RecordFailure()

		metrics.SMMProviderRequests.WithLabelValues(provider.Name, "status_error").Inc()

		logger.Logger.Error(
			"SMM provider status request failed",
			zap.Uint64("provider_id", provider.ID),
			zap.String("provider_name", provider.Name),
			zap.String("external_order_id", externalOrderID),
			zap.Error(gErr),
		)

		return "", gErr
	}

	breaker.RecordSuccess()

	metrics.SMMProviderRequests.
		WithLabelValues(provider.Name, "status_success").
		Inc()

	logger.Logger.Debug(
		"SMM provider status received",
		zap.Uint64("provider_id", provider.ID),
		zap.String("provider_name", provider.Name),
		zap.String("external_order_id", externalOrderID),
		zap.String("status", string(status)),
	)

	return status, nil
}
