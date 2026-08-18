package productservice

import (
	"context"
	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/productparams"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
	"time"

	"go.uber.org/zap"
)

//TODO - Later transferred to backoff service

func (s Service) CreateSMMMapping(ctx context.Context, req productparams.CreateSMMMappingRequest) (productparams.CreateSMMMappingResponse, error) {
	const Op = "productservice.AdminCreateSMMMapping"
	start := time.Now()

	smm := smmentity.SmmMapping{
		SmmServiceId: req.SmmServiceId,
		Name:         req.Name,
		Description:  req.Description,
		Category:     req.Category,
		Platform:     req.Platform,
		IsActive:     req.IsActive,
		ButtonName:   req.ButtonName,
	}

	if err := s.repository.SMMMappingCreate(ctx, &smm); err != nil {
		logger.Logger.Error("admin create smm mapping failed",
			zap.String("op", Op),
			zap.String("name", smm.Name),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return productparams.CreateSMMMappingResponse{}, richerror.New(Op, err).
			WithKind(richerror.KindInternal).
			WithMessage(msgerror.InternalServerError)
	}

	logger.Logger.Info("admin created smm mapping",
		zap.Int64("id", smm.Id),
		zap.String("name", smm.Name),
		zap.String("platform", string(req.Platform)),
		zap.String("category", string(req.Category)),
		zap.Duration("duration", time.Since(start)),
	)
	return productparams.CreateSMMMappingResponse{
		ID:           smm.Id,
		SmmServiceId: smm.SmmServiceId,
		Platform:     smm.Platform,
		Category:     smm.Category,
		ButtonName:   smm.ButtonName,
	}, nil
}
