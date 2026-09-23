package productservice

import (
	"context"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/productparams"
	"telegram-service-platform/pkg/metrics"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
	"time"

	"go.uber.org/zap"
)

func (s Service) GetSMMMappingByID(ctx context.Context, req productparams.GetSmmMappingByIDRequest) (productparams.GetSmmMappingByIDResponse, error) {
	const Op = "productservice.GetSMMMappingByID"
	start := time.Now()

	if s.smmCache != nil {
		if mapping, found, err := s.smmCache.GetMapping(ctx, req.Id); err == nil && found {
			metrics.SMMCacheHits.WithLabelValues("mapping").Inc()
			logger.Logger.Debug("smm mapping cache hit",
				zap.String("op", Op),
				zap.Int64("id", req.Id),
				zap.Duration("duration", time.Since(start)),
			)
			return productparams.GetSmmMappingByIDResponse{SmmMapping: mapping}, nil
		}
		metrics.SMMCacheMisses.WithLabelValues("mapping").Inc()
	}

	metrics.SMMCacheMisses.WithLabelValues("mapping").Inc()

	mapping, err := s.repository.SMMMappingGetByID(ctx, req.Id)
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

	if s.smmCache != nil {
		if setErr := s.smmCache.SetMapping(ctx, mapping); setErr != nil {
			logger.Logger.Warn(
				"failed to set smm mapping in cache",
				zap.String("op", Op),
				zap.Int64("id", req.Id),
				zap.Error(setErr),
			)
		}
	}

	logger.Logger.Debug("smm mapping found",
		zap.Int64("id", req.Id),
		zap.String("name", mapping.Name),
		zap.Duration("duration", time.Since(start)),
	)
	return productparams.GetSmmMappingByIDResponse{
		SmmMapping: mapping,
	}, nil
}
