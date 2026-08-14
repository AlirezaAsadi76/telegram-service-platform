package orderservice

import (
	"context"
	"telegram-service-platform/params/orderparams"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s *Service) GetByStatus(ctx context.Context, request orderparams.GetByStatusRequest) (orderparams.GetByStatusResponse, error) {
	const Op = "orderservice.getById"
	orders, iErr := s.repo.GetByStatus(ctx, request.Status)
	if iErr != nil {
		return orderparams.GetByStatusResponse{}, richerror.New(Op, iErr).WithKind(richerror.KindNotFound).WithMessage(msgerror.OrderNotFound)
	}
	return orderparams.GetByStatusResponse{
		Orders: orders,
	}, nil
}
