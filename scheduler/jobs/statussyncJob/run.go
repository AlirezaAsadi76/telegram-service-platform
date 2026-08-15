package statussyncjob

import (
	"context"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/notificationparams"
	"telegram-service-platform/params/orderparams"
	"telegram-service-platform/pkg/metrics"
	"time"

	"telegram-service-platform/entity/notificationentity"
	"telegram-service-platform/entity/orderentity"

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
	defer j.mutex.Unlock()

	resp, gErr := j.orderService.GetByStatus(ctx, orderparams.GetByStatusRequest{
		Status: orderentity.OrderStatusProcessing,
	})
	if gErr != nil {
		metrics.WorkerRuns.WithLabelValues(jobName, "error").Inc()
		logger.Logger.Error("worker failed", zap.String("job", jobName), zap.Error(gErr))
		return gErr
	}

	logger.Logger.Info("processing orders fetched",
		zap.String("job", jobName),
		zap.Int("count", len(resp.Orders)),
	)

	for _, order := range resp.Orders {
		logger.Logger.Debug("syncing order status",
			zap.String("job", jobName),
			zap.Uint64("order_id", order.ID),
			zap.String("external_order_id", order.ExternalOrderID),
		)

		if order.ProviderID == nil || order.ExternalOrderID == "" {
			continue
		}

		// نیاز به متد GetOrderStatus در smmproviderservice (فاز ۵)
		status, ggErr := j.smmProviderService.GetOrderStatus(ctx, *order.ProviderID, order.ExternalOrderID)
		if ggErr != nil {
			logger.Logger.Error("get order status failed", zap.String("job", jobName), zap.Error(ggErr))
			continue
		}

		switch status {
		case "COMPLETED":
			if err := j.orderService.UpdateStatus(ctx, orderparams.UpdateStatusRequest{OrderID: order.ID, Status: orderentity.OrderStatusSuccess}); err != nil {
				logger.Logger.Error("update order to completed failed", zap.String("job", jobName), zap.Error(err))
				continue
			}

			_ = j.notificationService.Create(ctx, notificationparams.CreateRequest{
				UserID: order.UserID,
				Type:   notificationentity.NotificationTypeOrderCompleted,
				Payload: map[string]any{
					"order_id": order.ID,
				},
			})

		case "FAILED", "CANCELLED":
			if err := j.orderService.UpdateStatus(ctx, orderparams.UpdateStatusRequest{
				OrderID: order.ID,
				Status:  orderentity.OrderStatusFailed,
			}); err != nil {
				logger.Logger.Error("update order to failed ", zap.String("job", jobName), zap.Error(err))
				continue
			}

			_ = j.notificationService.Create(ctx, notificationparams.CreateRequest{
				UserID: order.UserID,
				Type:   notificationentity.NotificationTypeOrderFailed,
				Payload: map[string]any{
					"order_id": order.ID,
					"reason":   "provider_status_failed",
				},
			})

		case "PROCESSING", "PENDING", "IN_PROGRESS":
			// No change
			logger.Logger.Info("No change - processing order status", zap.String("job", jobName))
		}
	}

	logger.Logger.Info("worker completed", zap.String("job", jobName), zap.Duration("duration", time.Since(start)))
	return nil
}
