package orderservice

import (
	"context"
	"telegram-service-platform/params/orderparams"
)

func (s *Service) UpdateOrderStatusById(ctx context.Context, req orderparams.UpdateOrderStatusByIdRequest) (*orderparams.UpdateOrderStatusByIdResponse, error) {

	resp := &orderparams.UpdateOrderStatusByIdResponse{}
	return resp, nil
}
