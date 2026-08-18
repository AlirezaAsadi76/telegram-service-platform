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

func (s Service) GetSMMMappingsByPlatformCategory(ctx context.Context, platform, category string) ([]smmentity.SmmMapping, error) {
	const Op = "productservice.GetSMMMappingsByPlatformCategory"
	start := time.Now()

	mappings, err := s.repository.SMMMappingGetByPlatformCategory(ctx, platform, category)
	if err != nil {
		logger.Logger.Error("get smm mappings by platform category failed",
			zap.String("op", Op),
			zap.String("platform", platform),
			zap.String("category", category),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return nil, richerror.New(Op, err).
			WithKind(richerror.KindInternal).
			WithMessage(msgerror.InternalServerError)
	}

	logger.Logger.Debug("smm mappings by platform category",
		zap.String("platform", platform),
		zap.String("category", category),
		zap.Int("count", len(mappings)),
		zap.Duration("duration", time.Since(start)),
	)
	return mappings, nil
}
