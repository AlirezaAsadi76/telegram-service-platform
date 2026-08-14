package statussyncjob

import (
	"context"
	"log"
	"telegram-service-platform/params/notificationparams"
	"telegram-service-platform/params/orderparams"

	"telegram-service-platform/entity/notificationentity"
	"telegram-service-platform/entity/orderentity"
)

func (j *Job) Run(ctx context.Context) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	resp, gErr := j.orderService.GetByStatus(ctx, orderparams.GetByStatusRequest{
		Status: orderentity.OrderStatusProcessing,
	})
	if gErr != nil {
		return gErr
	}

	for _, order := range resp.Orders {
		if order.ProviderID == nil || order.ExternalOrderID == "" {
			continue
		}

		// نیاز به متد GetOrderStatus در smmproviderservice (فاز ۵)
		status, gErr := j.smmProviderService.GetOrderStatus(ctx, *order.ProviderID, order.ExternalOrderID)
		if gErr != nil {
			log.Printf("get order status failed for order %d: %v", order.ID, gErr)
			continue
		}

		switch status {
		case "COMPLETED":
			if err := j.orderService.UpdateStatus(ctx, orderparams.UpdateStatusRequest{OrderID: order.ID, Status: orderentity.OrderStatusSuccess}); err != nil {
				log.Printf("update order to completed failed: %v", err)
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
				log.Printf("update order to failed failed: %v", err)
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
		}
	}

	return nil
}
