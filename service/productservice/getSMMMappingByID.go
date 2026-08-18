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

func (s Service) GetSMMMappingByID(ctx context.Context, req productparams.GetSmmMappingByIDRequest) (productparams.GetSmmMappingByIDResponse, error) {
	const Op = "productservice.GetSMMMappingByID"
	start := time.Now()

	m, err := s.repository.SMMMappingGetByID(ctx, req.Id)
	if err != nil {
		logger.Logger.Error("get smm mapping by id failed",
			zap.String("op", Op),
			zap.Int64("id", req.Id),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return productparams.GetSmmMappingByIDResponse{}, richerror.New(Op, err).
			WithKind(richerror.KindNotFound).
			WithMessage(msgerror.ProductNotFound)
	}

	logger.Logger.Debug("smm mapping found",
		zap.Int64("id", req.Id),
		zap.String("name", m.Name),
		zap.Duration("duration", time.Since(start)),
	)
	return productparams.GetSmmMappingByIDResponse{
		SmmMapping: m,
	}, nil
}
