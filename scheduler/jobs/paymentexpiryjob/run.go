package paymentexpiryjob

import (
	"context"
	"log"
	"telegram-service-platform/params/notificationparams"
	"telegram-service-platform/params/orderparams"
	"telegram-service-platform/params/paymentparams"

	"telegram-service-platform/entity/notificationentity"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/entity/paymententity"
)

func (j *Job) Run(ctx context.Context) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	resp, gErr := j.paymentService.GetExpired(ctx)
	if gErr != nil {
		return gErr
	}

	for _, payment := range resp.Payments {
		if err := j.paymentService.UpdateStatus(ctx, paymentparams.UpdateStatusRequest{
			PymentId: payment.ID,
			Status:   paymententity.PaymentStatusExpired,
		}); err != nil {
			log.Printf("update payment %d to expired failed: %v", payment.ID, err)
			continue
		}

		if err := j.orderService.UpdateStatus(ctx, orderparams.UpdateStatusRequest{
			OrderID: payment.ID,
			Status:  orderentity.OrderStatusCanceled,
		}); err != nil {
			log.Printf("cancel order %d failed: %v", payment.OrderID, err)
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

	return nil
}
