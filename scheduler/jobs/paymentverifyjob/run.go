package paymentverifyjob

import (
	"context"

	"telegram-service-platform/logger"
	"telegram-service-platform/params/notificationparams"

	"telegram-service-platform/params/paymentparams"
	"telegram-service-platform/pkg/metrics"
	"time"

	"telegram-service-platform/entity/notificationentity"

	"telegram-service-platform/entity/paymententity"

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

	j.mutex.Lock()
	defer j.mutex.Unlock()

	res, geErr := j.paymentService.GetPending(ctx)
	if geErr != nil {
		metrics.WorkerRuns.
			WithLabelValues(jobName, "error").
			Inc()

		logger.Logger.Error(
			"worker failed",
			zap.String("job", jobName),
			zap.Error(geErr),
		)

		return geErr
	}

	logger.Logger.Info(
		"worker pending payments fetched",
		zap.String("job", jobName),
		zap.Int("count", len(res.Payments)),
	)

	for _, payment := range res.Payments {
		logger.Logger.Debug(
			"verifying payment",
			zap.String("job", jobName),
			zap.Uint64("payment_id", payment.ID),
		)

		verifyResponse, err := j.paymentService.ConfirmPayment(
			ctx,
			paymentparams.ConfirmPaymentRequest{
				PaymentID: payment.ID,
			},
		)

		if err != nil {
			logger.Logger.Error(
				"verify failed",
				zap.String("job", jobName),
				zap.Uint64("payment_id", payment.ID),
				zap.Error(err),
			)

			updatedPayment, getErr := j.paymentService.GetById(
				ctx,
				payment.ID,
			)
			if getErr != nil {
				continue
			}

			if updatedPayment.Status == paymententity.PaymentStatusSuccess {
				logger.Logger.Warn(
					"payment succeeded but confirmation flow returned error; retrying confirmation for fulfillment enqueue",
					zap.Uint64("payment_id", updatedPayment.ID),
					zap.Uint64("order_id", updatedPayment.OrderID),
				)

				_, retryErr := j.paymentService.ConfirmPayment(
					ctx,
					paymentparams.ConfirmPaymentRequest{
						PaymentID: updatedPayment.ID,
					},
				)
				if retryErr != nil {
					logger.Logger.Error(
						"fulfillment enqueue recovery failed",
						zap.Uint64("payment_id", updatedPayment.ID),
						zap.Uint64("order_id", updatedPayment.OrderID),
						zap.Error(retryErr),
					)
				}
			}

			continue
		}

		switch verifyResponse.Status {
		case paymententity.PaymentStatusSuccess:
			_ = j.notificationService.Create(
				ctx,
				notificationparams.CreateRequest{
					UserID: payment.UserID,
					Type:   notificationentity.NotificationTypeOrderPaid,
					Payload: map[string]any{
						"order_id": payment.OrderID,
						"amount":   payment.Amount,
					},
				},
			)

		case paymententity.PaymentStatusFailed,
			paymententity.PaymentStatusCanceled,
			paymententity.PaymentStatusExpired:
			_ = j.notificationService.Create(
				ctx,
				notificationparams.CreateRequest{
					UserID: payment.UserID,
					Type:   notificationentity.NotificationTypeOrderFailed,
					Payload: map[string]any{
						"order_id": payment.OrderID,
						"reason":   "payment_failed",
					},
				},
			)

		case paymententity.PaymentStatusPending,
			paymententity.PaymentStatusUnknown:
			continue
		}
	}

	logger.Logger.Info(
		"worker completed",
		zap.String("job", jobName),
		zap.Duration("duration", time.Since(start)),
	)

	return nil
}
