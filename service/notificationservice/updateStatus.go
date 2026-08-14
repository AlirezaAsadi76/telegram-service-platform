package notificationservice

import (
	"context"
	"telegram-service-platform/params/notificationparams"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s *Service) UpdateStatus(ctx context.Context, request notificationparams.UpdateStatusRequest) error {
	const Op = "notificationservice.UpdateStatus"
	gErr := s.repo.UpdateStatus(ctx, request.Id, request.Status)
	if gErr != nil {
		return richerror.New(Op, gErr).WithKind(richerror.KindInternal).WithMessage(msgerror.InternalServerError)
	}

	return nil
}
