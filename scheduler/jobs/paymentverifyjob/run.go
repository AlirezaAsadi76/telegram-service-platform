package paymentverifyjob

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

	res, geErr := j.paymentService.GetPending(ctx)
	if geErr != nil {
		return geErr
	}

	for _, payment := range res.Payments {
		verifyResponse, err := j.paymentService.Verify(ctx, paymentparams.VerifyRequest{PaymentID: payment.ID})
		if err != nil {
			log.Printf("payment verify failed for payment %d: %v", payment.ID, err)
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
				log.Printf("update order status failed for order %d: %v", payment.OrderID, uErr)
				continue
			}

			// Push Order ID to Redis queue for OrderFulfillerJob
			// TODO- we need key insert in config
			if lErr := j.redis.LPush(ctx, j.config.QueueKey, payment.OrderID); lErr != nil {
				log.Printf("push to queue:orders:paid failed for order %d: %v", payment.OrderID, lErr)
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
				log.Printf("cancel order failed for order %d: %v", payment.OrderID, uErr)
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

	return nil
}
