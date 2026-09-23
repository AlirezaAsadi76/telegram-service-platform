package productservice

import (
	"context"
	"time"

	"telegram-service-platform/logger"
	"telegram-service-platform/pkg/metrics"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"

	"go.uber.org/zap"
)

func (s Service) SyncSMMServices(ctx context.Context) error {
	const op = "productservice.SyncSMMServices"

	start := time.Now()
	defer func() {
		metrics.SMMServiceSyncDuration.
			Observe(time.Since(start).Seconds())
	}()

	response, err := s.adapter.AllServices(ctx)
	if err != nil {
		metrics.SMMServiceSyncTotal.
			WithLabelValues("error").
			Inc()

		logger.Logger.Error(
			"smm service sync failed at provider",
			zap.String("op", op),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)

		return richerror.New(op, err).
			WithKind(richerror.KindExternalAPI).
			WithMessage(msgerror.ExternalServiceFailed)
	}

	if len(response.Services) == 0 {
		metrics.SMMServiceSyncTotal.
			WithLabelValues("empty").
			Inc()

		logger.Logger.Warn(
			"smm provider returned an empty service catalog",
			zap.String("op", op),
			zap.Duration("duration", time.Since(start)),
		)

		return nil
	}

	for _, service := range response.Services {
		if err := s.repository.SMMServiceCreateOrUpdate(
			ctx,
			service,
		); err != nil {
			metrics.SMMServiceSyncTotal.
				WithLabelValues("error").
				Inc()

			logger.Logger.Error(
				"failed to persist smm service",
				zap.String("op", op),
				zap.Int64("service_id", service.Service),
				zap.String("provider", service.ProviderName),
				zap.Error(err),
				zap.Duration("duration", time.Since(start)),
			)

			return richerror.New(op, err).
				WithKind(richerror.KindQueryFailure).
				WithMessage(msgerror.QueryFailed)
		}

		metrics.SMMServiceSyncTotal.
			WithLabelValues("success").
			Inc()

		if s.smmCache != nil {
			if err := s.smmCache.SetService(
				ctx,
				&service,
			); err != nil {
				logger.Logger.Warn(
					"failed to cache synced smm service",
					zap.String("op", op),
					zap.Int64("service_id", service.Service),
					zap.Error(err),
				)
			}
		}
	}

	logger.Logger.Info(
		"smm service sync completed",
		zap.String("op", op),
		zap.Int("count", len(response.Services)),
		zap.Duration("duration", time.Since(start)),
	)

	return nil
}
