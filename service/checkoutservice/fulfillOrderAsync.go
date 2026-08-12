package checkoutservice

import (
	"context"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/params/orderparams"
)

func (s *Service) fulfillOrderAsync(order *orderentity.Order) {
	ctx := context.Background()
	if err := s.smmSvc.FulfillOrder(ctx, order); err != nil {
		// Log error, will be retried by StatusSyncJob
		return
	}
	// Update order status
	s.orderSvc.UpdateStatus(ctx, orderparams.UpdateStatusRequest{
		OrderID:         order.ID,
		Status:          orderentity.OrderStatusProcessing,
		ExternalOrderID: order.ExternalOrderID,
	})
}
