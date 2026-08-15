package paymentexpiryjob

import (
	"context"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/notificationparams"
	"telegram-service-platform/params/orderparams"
	"telegram-service-platform/params/paymentparams"
	"telegram-service-platform/pkg/metrics"
	"time"

	"telegram-service-platform/entity/notificationentity"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/entity/paymententity"

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

	resp, gErr := j.paymentService.GetExpired(ctx)
	if gErr != nil {
		metrics.WorkerRuns.WithLabelValues(jobName, "error").Inc()
		logger.Logger.Error("worker failed", zap.String("job", jobName), zap.Error(gErr))
		return gErr
	}

	logger.Logger.Info("expired payments fetched",
		zap.String("job", jobName),
		zap.Int("count", len(resp.Payments)),
	)

	for _, payment := range resp.Payments {
		if err := j.paymentService.UpdateStatus(ctx, paymentparams.UpdateStatusRequest{
			PaymentId: payment.ID,
			Status:    paymententity.PaymentStatusExpired,
		}); err != nil {
			logger.Logger.Error("Update(EXPIRED) Status Payment failed",
				zap.String("job", jobName), zap.Uint64("payment_id", payment.ID),
				zap.Error(err))

			continue
		}

		if err := j.orderService.UpdateStatus(ctx, orderparams.UpdateStatusRequest{
			OrderID: payment.OrderID,
			Status:  orderentity.OrderStatusCanceled,
		}); err != nil {
			logger.Logger.Error("Update(CANCELED) Status Payment failed",
				zap.String("job", jobName), zap.Uint64("payment_id", payment.ID),
				zap.Error(err))
			continue
		}

		_ = j.notificationService.Create(ctx, notificationparams.CreateRequest{
			UserID: payment.UserID,
			Type:   notificationentity.NotificationTypePaymentExpired,
			Payload: map[string]any{
				"order_id": payment.OrderID,
				"amount":   payment.Amount,
			},
		})
	}

	logger.Logger.Info("worker completed", zap.String("job", jobName), zap.Duration("duration", time.Since(start)))

	return nil
}
