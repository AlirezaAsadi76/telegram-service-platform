package orderservice

import (
	"context"
	"telegram-service-platform/params/orderparams"
)

func (s *Service) CreateOrder(ctx context.Context, req orderparams.CreateOrderRequest) (*orderparams.CreateOrderResponse, error) {

	// check product exist
	// check amount correct
	// check access level user

	resp := &orderparams.CreateOrderResponse{}
	return resp, nil
}
