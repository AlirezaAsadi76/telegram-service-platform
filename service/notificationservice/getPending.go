package notificationservice

import (
	"context"
	"telegram-service-platform/params/notificationparams"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s *Service) GetPending(ctx context.Context, request notificationparams.GetPendingRequest) (notificationparams.GetPendingResponse, error) {
	const Op = "notificationservice.GetPending"

	notifications, gErr := s.repo.GetPending(ctx, request.Limit)
	if gErr != nil {
		return notificationparams.GetPendingResponse{}, richerror.New(Op, gErr).WithKind(richerror.KindInternal).WithMessage(msgerror.InternalServerError)
	}

	return notificationparams.GetPendingResponse{
		Notifications: notifications,
	}, nil
}
