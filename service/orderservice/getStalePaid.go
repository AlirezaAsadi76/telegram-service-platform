package orderservice

import (
	"context"
	"time"

	"telegram-service-platform/params/orderparams"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s *Service) GetStalePaid(ctx context.Context, olderThan time.Duration, limit int) (orderparams.GetByStatusResponse, error) {
	const Op = "orderservice.GetStalePaid"

	orders, err := s.repo.GetStalePaid(ctx, olderThan, limit)
	if err != nil {
		return orderparams.GetByStatusResponse{}, richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryFailed)
	}

	return orderparams.GetByStatusResponse{
		Orders: orders,
	}, nil
}
