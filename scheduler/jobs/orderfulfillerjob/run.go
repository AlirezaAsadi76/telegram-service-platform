package orderfulfillerjob

import (
	"context"
	"errors"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/notificationparams"
	"telegram-service-platform/pkg/metrics"
	"telegram-service-platform/pkg/unmarshal"
	"time"

	"telegram-service-platform/entity/notificationentity"
	"telegram-service-platform/entity/orderentity"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func (j *Job) Run(ctx context.Context) error {
	start := time.Now()
	jobName := j.Name()

	defer func() {
		metrics.WorkerDuration.
			WithLabelValues(jobName).
			Observe(time.Since(start).Seconds())
	}()

	logger.Logger.Info("worker started", zap.String("job", jobName))

	j.mutex.Lock()
	result, bErr := j.redis.BRPop(ctx, j.config.Timeout, j.config.QueueKey)
	j.mutex.Unlock()

	if bErr != nil {
		if errors.Is(bErr, redis.Nil) {
			logger.Logger.Debug("worker queue empty", zap.String("job", jobName))
			return nil
		}

		metrics.WorkerRuns.
			WithLabelValues(jobName, "error").
			Inc()

		logger.Logger.Error(
			"worker redis error",
			zap.String("job", jobName),
			zap.Error(bErr),
		)

		return bErr
	}

	orderID, uErr := unmarshal.UnmarshalToUint64(result[1])
	if uErr != nil {
		metrics.WorkerRuns.
			WithLabelValues(jobName, "error").
			Inc()

		logger.Logger.Error(
			"failed to unmarshal order id",
			zap.String("job", jobName),
			zap.Error(uErr),
		)

		return uErr
	}

	order, err := j.orderService.GetById(ctx, orderID)
	if err != nil {
		logger.Logger.Error(
			"get order failed",
			zap.String("job", jobName),
			zap.Uint64("order_id", orderID),
			zap.Error(err),
		)

		return nil
	}

	claimed, claimErr := j.orderService.ClaimForProcessing(ctx, order.ID)
	if claimErr != nil {
		metrics.WorkerRuns.
			WithLabelValues(jobName, "error").
			Inc()

		logger.Logger.Error(
			"failed to claim order for processing",
			zap.String("job", jobName),
			zap.Uint64("order_id", order.ID),
			zap.Error(claimErr),
		)

		return claimErr
	}

	if !claimed {
		logger.Logger.Info(
			"order was already claimed or is no longer paid",
			zap.String("job", jobName),
			zap.Uint64("order_id", order.ID),
			zap.String("status", string(order.Status)),
		)

		metrics.WorkerRuns.
			WithLabelValues(jobName, "already_claimed").
			Inc()

		return nil
	}

	logger.Logger.Info(
		"order claimed for fulfillment",
		zap.String("job", jobName),
		zap.Uint64("order_id", order.ID),
		zap.String("status", string(orderentity.OrderStatusProcessing)),
	)

	backoff := []time.Duration{
		10 * time.Second,
		30 * time.Second,
		60 * time.Second,
	}

	var lastErr error

	for attempt := 0; attempt < 4; attempt++ {
		if attempt > 0 {
			select {
			case <-time.After(backoff[attempt-1]):
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		logger.Logger.Info(
			"fulfill attempt",
			zap.String("job", jobName),
			zap.Uint64("order_id", order.ID),
			zap.Int("attempt", attempt+1),
		)

		fulErr := j.smmProviderService.FulfillOrder(ctx, order)

		if fulErr == nil {
			logger.Logger.Info(
				"fulfill succeeded",
				zap.String("job", jobName),
				zap.Uint64("order_id", order.ID),
			)

			metrics.SMMProviderRequests.
				WithLabelValues("default", "success").
				Inc()

			_ = j.notificationService.Create(
				ctx,
				notificationparams.CreateRequest{
					UserID: order.UserID,
					Type:   notificationentity.NotificationTypeOrderProcessing,
					Payload: map[string]any{
						"order_id": order.ID,
					},
				},
			)

			metrics.WorkerRuns.
				WithLabelValues(jobName, "success").
				Inc()

			return nil
		}

		lastErr = fulErr

		metrics.SMMProviderRequests.
			WithLabelValues("default", "error").
			Inc()

		logger.Logger.Warn(
			"fulfill failed",
			zap.String("job", jobName),
			zap.Uint64("order_id", order.ID),
			zap.Int("attempt", attempt+1),
			zap.Error(fulErr),
		)
	}

	logger.Logger.Error(
		"fulfill failed after retries",
		zap.String("job", jobName),
		zap.Uint64("order_id", order.ID),
		zap.Error(lastErr),
	)

	metrics.WorkerRuns.
		WithLabelValues(jobName, "failed_after_retries").
		Inc()

	// Do not refund or mark FAILED here yet.
	// The provider outcome can be unknown after a timeout.
	return nil
}
