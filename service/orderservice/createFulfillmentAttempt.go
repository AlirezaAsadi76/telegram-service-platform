package orderservice

import (
	"context"

	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s *Service) CreateFulfillmentAttempt(ctx context.Context, attempt *orderentity.FulfillmentAttempt) (uint64, error) {
	const Op = "orderservice.CreateFulfillmentAttempt"

	id, err := s.repo.CreateFulfillmentAttempt(ctx, attempt)
	if err != nil {
		return 0, richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.OrderUpdateFailed)
	}

	return id, nil
}
