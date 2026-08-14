package orderfulfillerjob

import (
	"context"
	"log"
	"telegram-service-platform/params/notificationparams"
	"telegram-service-platform/params/orderparams"
	"telegram-service-platform/pkg/unmarshal"
	"time"

	"telegram-service-platform/entity/notificationentity"
	"telegram-service-platform/entity/orderentity"
)

func (j *Job) Run(ctx context.Context) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	// BLPOP with 5s timeout — if empty, returns gracefully
	result, bErr := j.redis.BRPop(ctx, j.config.QueueKey, j.config.Timeout)
	if bErr != nil {
		return bErr
	}
	orderID, uErr := unmarshal.UnmarshalToUint64(result[1])
	if uErr != nil {
		return uErr
	}

	order, err := j.orderService.GetById(ctx, orderID)
	if err != nil {
		log.Printf("get order %d failed: %v", orderID, err)
		return nil
	}

	backoffs := []time.Duration{10 * time.Second, 30 * time.Second, 60 * time.Second}
	var lastErr error

	for attempt := 0; attempt <= 3; attempt++ {
		if attempt > 0 {
			time.Sleep(backoffs[attempt-1])
		}

		fulErr := j.smmProviderService.FulfillOrder(ctx, order)
		if fulErr == nil {
			if updateErr := j.orderService.UpdateStatus(ctx, orderparams.UpdateStatusRequest{
				OrderID: orderID,
				Status:  orderentity.OrderStatusProcessing,
			}); updateErr != nil {
				log.Printf("update order to processing failed: %v", updateErr)
			}
			_ = j.notificationService.Create(ctx, notificationparams.CreateRequest{
				UserID: order.UserID,
				Type:   notificationentity.NotificationTypeOrderPaid,
				Payload: map[string]any{
					"order_id": order.ID,
					"status":   "processing",
				}})
			return nil
		}

		lastErr = err
		log.Printf("fulfill order %d attempt %d failed: %v", order.ID, attempt+1, err)
	}

	log.Printf("fulfill order %d failed after 3 retries: %v", order.ID, lastErr)
	if uErr := j.orderService.UpdateStatus(ctx, orderparams.UpdateStatusRequest{
		OrderID:         order.ID,
		Status:          orderentity.OrderStatusFailed,
		ExternalOrderID: order.ExternalOrderID,
		ProviderID:      order.ProviderID,
	}); uErr != nil {
		log.Printf("update order to failed failed: %v", uErr)
	}

	_ = j.notificationService.Create(ctx, notificationparams.CreateRequest{
		UserID: order.UserID,
		Type:   notificationentity.NotificationTypeOrderFailed,
		Payload: map[string]any{
			"order_id": order.ID,
			"reason":   "provider_fulfill_failed",
		}})

	return nil
}
