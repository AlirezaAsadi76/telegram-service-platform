package paymentverifyjob

import (
	"context"
	"fmt"
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

	res, geErr := j.paymentService.GetPending(ctx)
	if geErr != nil {
		metrics.WorkerRuns.WithLabelValues(jobName, "error").Inc()
		logger.Logger.Error("worker failed", zap.String("job", jobName), zap.Error(geErr))
		return geErr
	}
	logger.Logger.Info("worker pending payments fetched",
		zap.String("job", jobName),
		zap.Int("count", len(res.Payments)),
	)

	for _, payment := range res.Payments {
		logger.Logger.Debug("verifying payment",
			zap.String("job", jobName),
			zap.Uint64("payment_id", payment.ID),
		)

		verifyResponse, err := j.paymentService.Verify(ctx, paymentparams.VerifyRequest{PaymentID: payment.ID})
		if err != nil {
			logger.Logger.Error("verify failed",
				zap.String("job", jobName),
				zap.Uint64("payment_id", payment.ID),
				zap.Error(err))

			continue
		}

		switch verifyResponse.Status {
		case paymententity.PaymentStatusSuccess:
			//TODO- we need Provider Id
			if uErr := j.orderService.UpdateStatus(ctx, orderparams.UpdateStatusRequest{
				Status:          orderentity.OrderStatusPaid,
				OrderID:         payment.OrderID,
				ExternalOrderID: payment.ExternalID,
			}); uErr != nil {
				logger.Logger.Error("update order failed",
					zap.String("job", jobName),
					zap.Uint64("payment_id", payment.ID),
					zap.Error(uErr))

				continue
			}

			// Push Order ID to Redis queue for OrderFulfillerJob
			// TODO- we need key insert in config
			if lErr := j.redis.LPush(ctx, j.config.QueueKey, payment.OrderID); lErr != nil {
				logger.Logger.Error("push to queue failed",
					zap.String("queue", j.config.QueueKey),
					zap.Uint64("order_id", payment.OrderID),
					zap.Error(lErr),
				)

			}

			_ = j.notificationService.Create(ctx, notificationparams.CreateRequest{
				UserID: payment.UserID,
				Type:   notificationentity.NotificationTypeOrderPaid,
				Payload: map[string]any{
					"order_id": payment.OrderID,
					"amount":   payment.Amount,
				},
			})

		case paymententity.PaymentStatusFailed, paymententity.PaymentStatusCanceled:
			if uErr := j.orderService.UpdateStatus(ctx, orderparams.UpdateStatusRequest{
				Status:  orderentity.OrderStatusCanceled,
				OrderID: payment.OrderID,
			}); uErr != nil {
				logger.Logger.Error(fmt.Sprintf("cancel order failed for order %d", payment.OrderID),
					zap.String("job", jobName),
					zap.Uint64("payment_id", payment.ID),
					zap.Error(uErr))

			}
			_ = j.notificationService.Create(ctx, notificationparams.CreateRequest{
				UserID: payment.UserID,
				Type:   notificationentity.NotificationTypeOrderFailed,
				Payload: map[string]any{
					"order_id": payment.OrderID,
					"reason":   "payment_failed",
				},
			})

		case paymententity.PaymentStatusPending:
			// Still pending, do nothing
		}
	}

	metrics.WorkerRuns.WithLabelValues(jobName, "success").Inc()
	logger.Logger.Info("worker completed", zap.String("job", jobName), zap.Duration("duration", time.Since(start)))
	return nil
}
