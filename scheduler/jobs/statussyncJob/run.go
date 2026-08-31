package statussyncjob

import (
	"context"
	"fmt"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/notificationparams"
	"telegram-service-platform/params/orderparams"
	"telegram-service-platform/params/walletparam"
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
			// ❌ اصلاح حیاتی: بازپرداخت خودکار وقتی ارائه‌دهنده سفارش را لغو یا رد می‌کند
			refundID := fmt.Sprintf("refund:status_cancelled:order:%d", order.ID)
			_, refErr := j.walletService.Credit(ctx, walletparam.CreditRequest{
				UserID:         order.UserID,
				Amount:         order.Amount,
				ReferenceID:    fmt.Sprintf("order:%d", order.ID),
				IdempotencyKey: refundID,
			})

			if refErr != nil {
				logger.Logger.Error("CRITICAL: AUTO-REFUND FAILED ON PROVIDER CANCELLATION",
					zap.Uint64("order_id", order.ID),
					zap.Error(refErr),
				)
				continue // وضعیت را تغییر نمی‌دهیم تا ادمین دستی بررسی کند
			}

			if err := j.orderService.UpdateStatus(ctx, orderparams.UpdateStatusRequest{
				OrderID: order.ID,
				Status:  orderentity.OrderStatusFailed, // یا OrderStatusCanceled اگر دارید
			}); err != nil {
				logger.Logger.Error("update order to failed/canceled failed", zap.String("job", jobName), zap.Error(err))
				continue
			}

			_ = j.notificationService.Create(ctx, notificationparams.CreateRequest{
				UserID: order.UserID,
				Type:   notificationentity.NotificationTypeOrderFailed,
				Payload: map[string]any{
					"order_id": order.ID,
					"reason":   "provider_cancelled_and_refunded",
				},
			})
			logger.Logger.Info("order cancelled by provider and user refunded", zap.Uint64("order_id", order.ID))
		case "PROCESSING", "PENDING", "IN_PROGRESS":
			// No change
			logger.Logger.Info("No change - processing order status", zap.String("job", jobName))
		}
	}

	logger.Logger.Info("worker completed", zap.String("job", jobName), zap.Duration("duration", time.Since(start)))
	return nil
}
