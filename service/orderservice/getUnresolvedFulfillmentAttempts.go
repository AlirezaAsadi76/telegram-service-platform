package orderservice

import (
	"context"

	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s *Service) GetUnresolvedFulfillmentAttempts(ctx context.Context, limit int) ([]orderentity.FulfillmentAttempt, error) {
	const Op = "orderservice.GetUnresolvedFulfillmentAttempts"

	attempts, err := s.repo.GetUnresolvedFulfillmentAttempts(ctx, limit)
	if err != nil {
		return nil, richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryFailed)
	}

	return attempts, nil
}
