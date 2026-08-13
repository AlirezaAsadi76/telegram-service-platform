package orderservice

import (
	"context"
	"telegram-service-platform/params/orderparams"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s *Service) UpdateStatus(ctx context.Context, req orderparams.UpdateStatusRequest) error {
	const Op = "orderservice.UpdateStatus"
	if err := s.repo.UpdateStatus(ctx, req.OrderID, req.Status, req.ExternalOrderID, req.ProviderID); err != nil {
		return richerror.New(Op, err).WithKind(richerror.KindQueryFailure).WithMessage(msgerror.OrderUpdateFailed)
	}
	return nil
}
