package productservice

import (
	"context"
	"time"

	"telegram-service-platform/logger"
	"telegram-service-platform/params/productparams"
	"telegram-service-platform/pkg/metrics"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"

	"go.uber.org/zap"
)

func (s Service) GetSMMServiceByID(ctx context.Context, req productparams.GetSmmServiceByIDRequest) (productparams.GetSmmServiceByIDResponse, error) {
	const op = "productservice.GetSMMServiceByID"
	start := time.Now()
	if s.smmCache != nil {
		if service, found, err := s.smmCache.GetService(ctx, req.Id); err == nil && found {
			metrics.SMMCacheHits.WithLabelValues("service").Inc()
			logger.Logger.Debug("smm service cache hit",
				zap.String("op", op),
				zap.Int64("id", req.Id),
				zap.Duration("duration", time.Since(start)),
			)
			return productparams.GetSmmServiceByIDResponse{Smm: service}, nil
		}

		metrics.SMMCacheMisses.WithLabelValues("service").Inc()
	}

	service, err := s.repository.SMMServiceGetByD(ctx, req.Id)
	if err != nil {
		logger.Logger.Error("get smm service by id failed (db)",
			zap.String("op", op),
			zap.Int64("id", req.Id),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return productparams.GetSmmServiceByIDResponse{}, richerror.New(op, err).
			WithKind(richerror.KindNotFound).
			WithMessage(msgerror.ProductNotFound)
	}

	if s.smmCache != nil {
		if setErr := s.smmCache.SetService(ctx, service); setErr != nil {
			logger.Logger.Warn(
				"failed to set smm service in cache",
				zap.String("op", op),
				zap.Int64("id", req.Id),
				zap.Error(setErr),
			)
		}
	}

	logger.Logger.Debug("smm service found and cached",
		zap.Int64("id", req.Id),
		zap.String("name", service.Name),
		zap.Duration("duration", time.Since(start)),
	)

	return productparams.GetSmmServiceByIDResponse{
		Smm: service,
	}, nil
}
