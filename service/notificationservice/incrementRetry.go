package notificationservice

import (
	"context"
	"telegram-service-platform/params/notificationparams"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s *Service) IncrementRetry(ctx context.Context, request notificationparams.IncrementRetryRequest) error {
	const Op = "notificationservice.IncrementRetry"

	gErr := s.repo.UpdateRetryCount(ctx, request.Id, request.RetryCount)
	if gErr != nil {
		return richerror.New(Op, gErr).WithKind(richerror.KindInternal).WithMessage(msgerror.InternalServerError)
	}

	return nil

}
