// service/orderservice/create.go
package orderservice

import (
	"context"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/params/orderparams"
	"telegram-service-platform/pkg/richerror"
)

func (s *Service) Create(ctx context.Context, req orderparams.CreateRequest) (*orderparams.CreateResponse, error) {
	const Op = "OrderService.Create"

	order := &orderentity.Order{
		UserID:      req.UserID,
		ProductType: req.ProductType,
		ProductID:   req.ProductID,
		Quantity:    req.Quantity,
		TargetLink:  req.TargetLink,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Status:      orderentity.OrderStatusPending,
	}
	if err := s.repo.Create(ctx, order); err != nil {
		return nil, richerror.New(Op, err).WithKind(richerror.KindCreateFailed)
	}
	return &orderparams.CreateResponse{OrderID: order.ID}, nil
}
