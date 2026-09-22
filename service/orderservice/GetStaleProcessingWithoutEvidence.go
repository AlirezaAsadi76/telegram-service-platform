package orderservice

import (
	"context"
	"time"

	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s *Service) GetStaleProcessingWithoutEvidence(ctx context.Context, olderThan time.Duration, limit int) ([]*orderentity.Order, error) {
	const Op = "orderservice.GetStaleProcessingWithoutEvidence"

	orders, err := s.repo.GetStaleProcessingWithoutEvidence(
		ctx,
		olderThan,
		limit,
	)
	if err != nil {
		return nil, richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryFailed)
	}

	return orders, nil
}
