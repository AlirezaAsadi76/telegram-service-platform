package orderservice

import (
	"context"

	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s *Service) CompleteProcessing(ctx context.Context, orderID uint64) (bool, error) {
	const Op = "orderservice.CompleteProcessing"

	completed, err := s.repo.CompleteProcessing(ctx, orderID)
	if err != nil {
		return false, richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.OrderUpdateFailed)
	}

	return completed, nil
}
