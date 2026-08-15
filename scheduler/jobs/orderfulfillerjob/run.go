package orderfulfillerjob

import (
	"context"
	"errors"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/notificationparams"
	"telegram-service-platform/params/orderparams"
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
		metrics.WorkerDuration.WithLabelValues(jobName).Observe(time.Since(start).Seconds())
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
		metrics.WorkerRuns.WithLabelValues(jobName, "error").Inc()
		logger.Logger.Error("worker redis error", zap.String("job", jobName), zap.Error(bErr))
		return bErr
	}
	orderID, uErr := unmarshal.UnmarshalToUint64(result[1])
	if uErr != nil {
		return uErr
	}

	order, err := j.orderService.GetById(ctx, orderID)
	if err != nil {
		logger.Logger.Error("get order failed", zap.Uint64("order_id", orderID), zap.Error(err))
		return nil
	}

	backoffs := []time.Duration{10 * time.Second, 30 * time.Second, 60 * time.Second}
	var lastErr error

	for attempt := 0; attempt <= 3; attempt++ {
		logger.Logger.Info("fulfill attempt",
			zap.String("job", jobName),
			zap.Uint64("order_id", order.ID),
			zap.Int("attempt", attempt+1),
		)

		if attempt > 0 {
			time.Sleep(backoffs[attempt-1])
		}

		fulErr := j.smmProviderService.FulfillOrder(ctx, order)
		if fulErr == nil {
			if updateErr := j.orderService.UpdateStatus(ctx, orderparams.UpdateStatusRequest{
				OrderID: orderID,
				Status:  orderentity.OrderStatusProcessing,
			}); updateErr != nil {
				logger.Logger.Error("update order status failed",
					zap.String("job", jobName),
					zap.Uint64("order_id", order.ID),
					zap.Error(updateErr))

			}

			metrics.SMMProviderRequests.WithLabelValues("default", "success").Inc()
			metrics.WorkerRuns.WithLabelValues(jobName, "success").Inc()
			logger.Logger.Info("fulfill succeeded",
				zap.String("job", jobName),
				zap.Uint64("order_id", order.ID),
			)

			_ = j.notificationService.Create(ctx, notificationparams.CreateRequest{
				UserID: order.UserID,
				Type:   notificationentity.NotificationTypeOrderPaid,
				Payload: map[string]any{
					"order_id": order.ID,
					"status":   "processing",
				}})
			return nil
		}
		metrics.SMMProviderRequests.WithLabelValues("default", "error").Inc()
		logger.Logger.Warn("fulfill failed",
			zap.String("job", jobName),
			zap.Uint64("order_id", order.ID),
			zap.Int("attempt", attempt+1),
			zap.Error(fulErr),
		)
		lastErr = fulErr

	}

	if upErr := j.orderService.UpdateStatus(ctx, orderparams.UpdateStatusRequest{
		OrderID:         order.ID,
		Status:          orderentity.OrderStatusFailed,
		ExternalOrderID: order.ExternalOrderID,
		ProviderID:      order.ProviderID,
	}); upErr != nil {
		logger.Logger.Error("update order status failed",
			zap.String("job", jobName),
			zap.Uint64("order_id", order.ID),
			zap.Error(upErr))
	}

	_ = j.notificationService.Create(ctx, notificationparams.CreateRequest{
		UserID: order.UserID,
		Type:   notificationentity.NotificationTypeOrderFailed,
		Payload: map[string]any{
			"order_id": order.ID,
			"reason":   "provider_fulfill_failed",
		}})

	metrics.WorkerRuns.WithLabelValues(jobName, "failed_after_retries").Inc()
	logger.Logger.Error("fulfill failed after retries",
		zap.String("job", jobName),
		zap.Uint64("order_id", order.ID),
		zap.Error(lastErr),
	)

	return nil
}
