package orderservice

import (
	"context"
	"telegram-service-platform/params/orderparams"
)

func (s *Service) GetOrderById(ctx context.Context, req orderparams.GetOrderByIdRequest) (*orderparams.GetOrderByIdResponse, error) {

	resp := &orderparams.GetOrderByIdResponse{}
	return resp, nil
}
