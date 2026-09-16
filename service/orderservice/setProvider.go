package orderservice

import (
	"context"

	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s *Service) AssignProvider(ctx context.Context, orderID uint64, providerID uint64) error {
	const Op = "orderservice.AssignProvider"

	if err := s.repo.AssignProvider(ctx, orderID, providerID); err != nil {
		return richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.OrderUpdateFailed)
	}

	return nil
}
