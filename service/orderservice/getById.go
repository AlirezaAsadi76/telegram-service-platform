package orderservice

import (
	"context"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s *Service) GetById(ctx context.Context, orderId uint64) (*orderentity.Order, error) {
	const Op = "orderservice.getById"
	order, iErr := s.repo.GetByID(ctx, orderId)
	if iErr != nil {
		return nil, richerror.New(Op, iErr).WithKind(richerror.KindNotFound).WithMessage(msgerror.OrderNotFound)
	}
	return order, nil
}
