package orderservice

import (
	"context"

	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s *Service) MarkFulfillmentAttemptResolved(ctx context.Context, attemptID uint64) error {
	const Op = "orderservice.MarkFulfillmentAttemptResolved"

	if err := s.repo.MarkFulfillmentAttemptResolved(ctx, attemptID); err != nil {
		return richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.OrderUpdateFailed)
	}

	return nil
}
