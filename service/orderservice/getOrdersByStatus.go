package orderservice

import (
	"context"
	"telegram-service-platform/params/orderparams"
)

func (s *Service) GetOrdersByStatus(ctx context.Context, req orderparams.GetOrdersByStatusRequest) (*orderparams.GetOrdersByStatusResponse, error) {

	resp := &orderparams.GetOrdersByStatusResponse{}
	return resp, nil
}
