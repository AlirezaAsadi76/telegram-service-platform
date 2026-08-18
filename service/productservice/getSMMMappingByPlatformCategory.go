package productservice

import (
	"context"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/productparams"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
	"time"

	"go.uber.org/zap"
)

func (s Service) GetSMMMappingsByPlatformCategory(ctx context.Context, req productparams.GetSmmMappingByPlatformCategoryRequest) (
	productparams.GetSmmMappingByPlatformCategoryResponse, error) {
	const Op = "productservice.GetSMMMappingsByPlatformCategory"
	start := time.Now()

	mappings, err := s.repository.SMMMappingGetByPlatformCategory(ctx, req.Platform, req.Category)
	if err != nil {
		logger.Logger.Error("get smm mappings by platform category failed",
			zap.String("op", Op),
			zap.String("platform", string(req.Platform)),
			zap.String("category", string(req.Category)),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return productparams.GetSmmMappingByPlatformCategoryResponse{}, richerror.New(Op, err).
			WithKind(richerror.KindInternal).
			WithMessage(msgerror.InternalServerError)
	}

	logger.Logger.Debug("smm mappings by platform category",
		zap.String("platform", string(req.Platform)),
		zap.String("category", string(req.Category)),
		zap.Int("count", len(mappings)),
		zap.Duration("duration", time.Since(start)),
	)
	return productparams.GetSmmMappingByPlatformCategoryResponse{
		SmmMapping: mappings,
	}, nil
}
