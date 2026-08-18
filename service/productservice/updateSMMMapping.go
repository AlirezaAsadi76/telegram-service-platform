package productservice

import (
	"context"
	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/logger"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
	"time"

	"go.uber.org/zap"
)

//TODO - Later transferred to backoff service

func (s Service) UpdateSMMMapping(ctx context.Context, m *smmentity.SmmMapping) error {
	const Op = "productservice.AdminUpdateSMMMapping"
	start := time.Now()

	if err := s.repository.SMMMappingUpdate(ctx, m); err != nil {
		logger.Logger.Error("admin update smm mapping failed",
			zap.String("op", Op),
			zap.Int64("id", m.Id),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return richerror.New(Op, err).
			WithKind(richerror.KindInternal).
			WithMessage(msgerror.InternalServerError)
	}

	logger.Logger.Info("admin updated smm mapping",
		zap.Int64("id", m.Id),
		zap.String("name", m.Name),
		zap.Duration("duration", time.Since(start)),
	)
	return nil
}
