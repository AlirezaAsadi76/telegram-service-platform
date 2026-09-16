package orderservice

import (
	"context"

	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s *Service) SetProviderOrder(ctx context.Context, orderID uint64, providerID uint64, externalOrderID string) error {
	const Op = "orderservice.SetProviderOrder"

	if err := s.repo.SetProviderOrder(ctx, orderID, providerID, externalOrderID); err != nil {
		return richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.OrderUpdateFailed)
	}

	return nil
}
