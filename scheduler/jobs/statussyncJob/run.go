package statussyncjob

import (
	"context"
	"time"

	"telegram-service-platform/entity/notificationentity"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/notificationparams"
	"telegram-service-platform/params/orderparams"
	"telegram-service-platform/pkg/metrics"

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

	logger.Logger.Info(
		"status sync worker started",
		zap.String("job", jobName),
	)

	j.mutex.Lock()
	defer j.mutex.Unlock()

	resp, err := j.orderService.GetByStatus(
		ctx,
		orderparams.GetByStatusRequest{
			Status: orderentity.OrderStatusProcessing,
		},
	)
	if err != nil {
		metrics.WorkerRuns.
			WithLabelValues(jobName, "fetch_error").
			Inc()

		logger.Logger.Error(
			"status sync failed to fetch processing orders",
			zap.String("job", jobName),
			zap.Error(err),
		)

		return err
	}

	logger.Logger.Info(
		"processing orders fetched for status sync",
		zap.String("job", jobName),
		zap.Int("count", len(resp.Orders)),
	)

	for _, order := range resp.Orders {
		if order.ProviderID == nil || order.ExternalOrderID == "" {
			metrics.WorkerRuns.
				WithLabelValues(jobName, "missing_provider_data").
				Inc()

			logger.Logger.Warn(
				"processing order has incomplete provider data",
				zap.String("job", jobName),
				zap.Uint64("order_id", order.ID),
				zap.Bool("provider_id_missing", order.ProviderID == nil),
				zap.Bool("external_order_id_missing", order.ExternalOrderID == ""),
			)

			continue
		}

		logger.Logger.Debug(
			"syncing order status",
			zap.String("job", jobName),
			zap.Uint64("order_id", order.ID),
			zap.Uint64("provider_id", *order.ProviderID),
			zap.String("external_order_id", order.ExternalOrderID),
		)

		status, err := j.smmProviderService.GetOrderStatus(
			ctx,
			*order.ProviderID,
			order.ExternalOrderID,
		)
		if err != nil {
			metrics.WorkerRuns.
				WithLabelValues(jobName, "provider_status_error").
				Inc()

			logger.Logger.Error(
				"failed to get provider order status",
				zap.String("job", jobName),
				zap.Uint64("order_id", order.ID),
				zap.Uint64("provider_id", *order.ProviderID),
				zap.String("external_order_id", order.ExternalOrderID),
				zap.Error(err),
			)

			continue
		}

		switch status {
		case orderentity.OrderStatusCompleted:
			completed, err := j.orderService.CompleteProcessing(
				ctx,
				order.ID,
			)
			if err != nil {
				metrics.WorkerRuns.
					WithLabelValues(jobName, "complete_failed").
					Inc()

				logger.Logger.Error(
					"failed to complete processing order",
					zap.String("job", jobName),
					zap.Uint64("order_id", order.ID),
					zap.Error(err),
				)

				continue
			}

			if !completed {
				metrics.WorkerRuns.
					WithLabelValues(jobName, "order_state_changed").
					Inc()

				logger.Logger.Info(
					"order was no longer processing when completion was attempted",
					zap.String("job", jobName),
					zap.Uint64("order_id", order.ID),
				)

				continue
			}

			metrics.WorkerRuns.
				WithLabelValues(jobName, "completed").
				Inc()

			logger.Logger.Info(
				"order completed by provider status sync",
				zap.String("job", jobName),
				zap.Uint64("order_id", order.ID),
				zap.Uint64("provider_id", *order.ProviderID),
				zap.String("external_order_id", order.ExternalOrderID),
			)

			if err := j.notificationService.Create(
				ctx,
				notificationparams.CreateRequest{
					UserID: order.UserID,
					Type:   notificationentity.NotificationTypeOrderCompleted,
					Payload: map[string]any{
						"order_id": order.ID,
					},
				},
			); err != nil {
				logger.Logger.Error(
					"failed to create order completed notification",
					zap.String("job", jobName),
					zap.Uint64("order_id", order.ID),
					zap.Error(err),
				)
			}

		case orderentity.OrderStatusProcessing:
			logger.Logger.Debug(
				"provider order is still processing",
				zap.String("job", jobName),
				zap.Uint64("order_id", order.ID),
			)

		case orderentity.OrderStatusFailed:
			metrics.WorkerRuns.
				WithLabelValues(jobName, "provider_failed_pending_refund").
				Inc()

			logger.Logger.Info(
				"provider marked order failed; refund handled in next checkpoint",
				zap.String("job", jobName),
				zap.Uint64("order_id", order.ID),
			)

		default:
			metrics.WorkerRuns.
				WithLabelValues(jobName, "unsupported_status").
				Inc()

			logger.Logger.Warn(
				"unsupported provider status returned",
				zap.String("job", jobName),
				zap.Uint64("order_id", order.ID),
				zap.String("status", string(status)),
			)
		}
	}

	logger.Logger.Info(
		"status sync worker completed",
		zap.String("job", jobName),
		zap.Duration("duration", time.Since(start)),
	)

	return nil
}
