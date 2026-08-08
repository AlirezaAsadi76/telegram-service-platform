package orderservice

import (
	"context"
	"telegram-service-platform/params/orderparams"
)

func (s *Service) GetOrdersByUserId(ctx context.Context, req orderparams.GetOrdersByUserIdRequest) (*orderparams.GetOrdersByUserIdResponse, error) {

	resp := &orderparams.GetOrdersByUserIdResponse{}
	return resp, nil
}
