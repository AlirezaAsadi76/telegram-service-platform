// service/notificationservice/create.go
package notificationservice

import (
	"context"
	"telegram-service-platform/entity/notificationentity"
	"telegram-service-platform/params/notificationparams"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s *Service) Create(ctx context.Context, request notificationparams.CreateRequest) error {
	const Op = "notificationservice.Create"
	notification := &notificationentity.Notification{
		UserID:     request.UserID,
		Type:       request.Type,
		Status:     notificationentity.NotificationStatusPending,
		Payload:    request.Payload,
		RetryCount: 0,
	}
	if err := s.repo.Create(ctx, notification); err != nil {
		return richerror.New(Op, err).WithKind(richerror.KindInternal).WithMessage(msgerror.InternalServerError)
	}

	if err := s.redis.LPush(ctx, s.config.queueKey, notification); err != nil {
		return richerror.New(Op, err).WithKind(richerror.KindInternal).WithMessage(msgerror.InternalServerError)
	}

	return nil
}
