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

func (s Service) CreateSMMMapping(ctx context.Context, m *smmentity.SmmMapping) error {
	const Op = "productservice.AdminCreateSMMMapping"
	start := time.Now()

	if err := s.repository.SMMMappingCreate(ctx, m); err != nil {
		logger.Logger.Error("admin create smm mapping failed",
			zap.String("op", Op),
			zap.String("name", m.Name),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return richerror.New(Op, err).
			WithKind(richerror.KindInternal).
			WithMessage(msgerror.InternalServerError)
	}

	logger.Logger.Info("admin created smm mapping",
		zap.Int64("id", m.Id),
		zap.String("name", m.Name),
		zap.String("platform", m.Platform),
		zap.String("category", m.Category),
		zap.Duration("duration", time.Since(start)),
	)
	return nil
}
