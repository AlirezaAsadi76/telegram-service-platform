package productservice

import (
	"context"
	"time"

	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/logger"
	"telegram-service-platform/pkg/richerror"

	"go.uber.org/zap"
)

func (s Service) GetMissingSMMServices(ctx context.Context) ([]smmentity.SMM, error) {
	const op = "productservice.GetMissingSMMServices"

	start := time.Now()

	adapterServices, err := s.adapter.AllServices(ctx)
	if err != nil {
		logger.Logger.Error(
			"failed to fetch SMM services from provider",
			zap.String("op", op),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)

		return nil, richerror.New(op, err).
			WithKind(richerror.KindExternalAPI)
	}

	dbServices, err := s.repository.SMMServiceGetAll(ctx)
	if err != nil {
		logger.Logger.Error(
			"failed to fetch SMM services from database",
			zap.String("op", op),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)

		return nil, richerror.New(op, err).
			WithKind(richerror.KindQueryFailure)
	}

	adapterMap := make(map[int64]struct{}, len(adapterServices.Services))

	for _, service := range adapterServices.Services {
		adapterMap[service.Service] = struct{}{}
	}

	missingServices := make(
		[]smmentity.SMM,
		0,
	)

	for _, dbService := range dbServices {
		if _, exists := adapterMap[dbService.Service]; exists {
			continue
		}

		missingServices = append(
			missingServices,
			dbService,
		)
	}

	logger.Logger.Debug(
		"SMM missing-service validation completed",
		zap.String("op", op),
		zap.Int("provider_services", len(adapterServices.Services)),
		zap.Int("database_services", len(dbServices)),
		zap.Int("missing_services", len(missingServices)),
		zap.Duration("duration", time.Since(start)),
	)

	return missingServices, nil
}
