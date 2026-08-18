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

func (s Service) UpdateSMMMapping(ctx context.Context, req productparams.UpdateSMMMappingRequest) (productparams.UpdateSMMMappingResponse, error) {
	const Op = "productservice.AdminUpdateSMMMapping"
	start := time.Now()

	smm := smmentity.SmmMapping{
		Id:           req.Id,
		SmmServiceId: req.SmmServiceId,
		Name:         req.Name,
		Description:  req.Description,
		Category:     req.Category,
		Platform:     req.Platform,
		IsActive:     req.IsActive,
		ButtonName:   req.ButtonName,
	}
	if err := s.repository.SMMMappingUpdate(ctx, &smm); err != nil {
		logger.Logger.Error("admin update smm mapping failed",
			zap.String("op", Op),
			zap.Int64("id", smm.Id),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return productparams.UpdateSMMMappingResponse{}, richerror.New(Op, err).
			WithKind(richerror.KindInternal).
			WithMessage(msgerror.InternalServerError)
	}

	logger.Logger.Info("admin updated smm mapping",
		zap.Int64("id", smm.Id),
		zap.String("name", smm.Name),
		zap.Duration("duration", time.Since(start)),
	)
	return productparams.UpdateSMMMappingResponse{
		ID:           smm.Id,
		SmmServiceId: smm.SmmServiceId,
		Platform:     smm.Platform,
		Category:     smm.Category,
		ButtonName:   smm.ButtonName,
	}, nil
}
