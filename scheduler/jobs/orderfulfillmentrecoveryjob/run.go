package orderfulfillmentrecoveryjob

import (
	"context"
	"telegram-service-platform/logger"
	"telegram-service-platform/pkg/metrics"
	"telegram-service-platform/pkg/richerror"
	"time"

	"go.uber.org/zap"
)

func (j *Job) Run(ctx context.Context) error {
	const Op = "orderfulfillmentrecoveryjob.Run"

	start := time.Now()
	jobName := j.Name()

	defer func() {
		metrics.WorkerDuration.
			WithLabelValues(jobName).
			Observe(time.Since(start).Seconds())
	}()

	j.mutex.Lock()
	defer j.mutex.Unlock()

	if err := j.recoverUnresolvedAttempts(ctx); err != nil {
		metrics.WorkerRuns.
			WithLabelValues(jobName, "attempt_recovery_error").
			Inc()

		logger.Logger.Error(
			"order fulfillment attempt recovery failed",
			zap.String("job", jobName),
			zap.Error(err),
		)

		return richerror.New(Op, err)
	}

	result, err := j.orderService.GetStalePaid(
		ctx,
		j.config.StaleAfter,
		j.config.BatchSize,
	)
	if err != nil {
		metrics.WorkerRuns.
			WithLabelValues(jobName, "error").
			Inc()

		logger.Logger.Error(
			"order fulfillment recovery failed to fetch stale orders",
			zap.String("job", jobName),
			zap.Error(err),
		)

		return richerror.New(Op, err)
	}

	if len(result.Orders) == 0 {
		logger.Logger.Debug(
			"no stale paid orders found",
			zap.String("job", jobName),
		)

		metrics.WorkerRuns.
			WithLabelValues(jobName, "success").
			Inc()

		return nil
	}

	var failed int

	for _, order := range result.Orders {
		if err := j.fulfillmentEnqueuer.Enqueue(
			ctx,
			order.ID,
		); err != nil {
			failed++

			logger.Logger.Error(
				"failed to re-enqueue stale paid order",
				zap.String("job", jobName),
				zap.Uint64("order_id", order.ID),
				zap.Error(err),
			)

			continue
		}

		logger.Logger.Info(
			"stale paid order re-enqueued",
			zap.String("job", jobName),
			zap.Uint64("order_id", order.ID),
			zap.String("status", string(order.Status)),
		)
	}

	if err := j.recoverStaleProcessingWithoutEvidence(ctx); err != nil {
		metrics.WorkerRuns.
			WithLabelValues(jobName, "processing_recovery_error").
			Inc()

		logger.Logger.Error(
			"order fulfillment recovery failed to inspect stale processing orders",
			zap.String("job", jobName),
			zap.Error(err),
		)

		return richerror.New(Op, err)
	}

	if len(result.Orders) == 0 && failed == 0 {
		logger.Logger.Debug(
			"no stale paid orders found",
			zap.String("job", jobName),
		)
	}

	if failed > 0 {
		metrics.WorkerRuns.
			WithLabelValues(jobName, "partial_error").
			Inc()
	} else {
		metrics.WorkerRuns.
			WithLabelValues(jobName, "success").
			Inc()
	}

	logger.Logger.Info(
		"order fulfillment recovery completed",
		zap.String("job", jobName),
		zap.Int("stale_paid_orders", len(result.Orders)),
		zap.Int("reenqueue_failed", failed),
		zap.Duration("duration", time.Since(start)),
	)

	return nil
}
