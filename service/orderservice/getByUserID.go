package orderservice

import (
	"context"
	"telegram-service-platform/params/orderparams"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s *Service) GetByUserID(ctx context.Context, req orderparams.GetByUserIdRequest) (orderparams.GetByUserIdResponse, error) {
	const op = "orderservice.GetByUserID"

	orders, err := s.repo.GetByUserID(ctx, req.UserID)
	if err != nil {
		return orderparams.GetByUserIdResponse{}, richerror.New(op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryFailed)
	}

	return orderparams.GetByUserIdResponse{Orders: orders}, nil
}
