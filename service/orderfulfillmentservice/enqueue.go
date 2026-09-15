package orderfulfillmentservice

import (
	"context"
	"time"

	"telegram-service-platform/logger"
	"telegram-service-platform/pkg/metrics"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"

	"go.uber.org/zap"
)

func (s *Service) Enqueue(ctx context.Context, orderID uint64) error {
	const Op = "orderfulfillmentservice.Enqueue"

	start := time.Now()
	defer func() {
		metrics.OrderFulfillmentEnqueueDuration.Observe(
			time.Since(start).Seconds(),
		)
	}()

	if err := s.queue.LPush(ctx, s.config.QueueKey, orderID); err != nil {
		metrics.OrderFulfillmentEnqueueTotal.
			WithLabelValues("error").
			Inc()

		logger.Logger.Error(
			"failed to enqueue order fulfillment",
			zap.Uint64("order_id", orderID),
			zap.String("queue", s.config.QueueKey),
			zap.Error(err),
		)

		return richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryFailed)
	}

	metrics.OrderFulfillmentEnqueueTotal.
		WithLabelValues("success").
		Inc()

	logger.Logger.Info(
		"order fulfillment enqueued",
		zap.Uint64("order_id", orderID),
		zap.String("queue", s.config.QueueKey),
	)

	return nil
}
