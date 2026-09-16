package orderservice

import (
	"context"

	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s *Service) SaveExternalOrder(ctx context.Context, orderID uint64, providerID uint64, externalOrderID string) error {
	const Op = "orderservice.SaveExternalOrder"

	if err := s.repo.SaveExternalOrder(ctx, orderID, providerID, externalOrderID); err != nil {
		return richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.OrderUpdateFailed)
	}

	return nil
}
