package productservice

import (
	"context"
	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/logger"
	"telegram-service-platform/pkg/metrics"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
	"time"

	"go.uber.org/zap"
)

// Use this to build Telegram bot inline keyboards (first level: platforms, then categories, then services).
func (s Service) GetSMMCatalog(ctx context.Context) (map[string]map[string][]smmentity.SmmMapping, error) {
	const Op = "productservice.GetSMMCatalog"
	start := time.Now()

	mappings, rErr := s.repository.SMMMappingGetAllActive(ctx)
	if rErr != nil {
		logger.Logger.Error("get smm catalog failed",
			zap.String("op", Op),
			zap.Error(rErr),
			zap.Duration("duration", time.Since(start)),
		)
		return nil, richerror.New(Op, rErr).
			WithKind(richerror.KindInternal).
			WithMessage(msgerror.InternalServerError)
	}

	catalog := make(map[string]map[string][]smmentity.SmmMapping)
	for _, m := range mappings {
		if catalog[m.Platform] == nil {
			catalog[m.Platform] = make(map[string][]smmentity.SmmMapping)
		}
		catalog[m.Platform][m.Category] = append(catalog[m.Platform][m.Category], m)
	}

	metrics.SMMCatalogSize.Set(float64(len(mappings)))
	logger.Logger.Debug("smm catalog built",
		zap.Int("platforms", len(catalog)),
		zap.Int("mappings", len(mappings)),
		zap.Duration("duration", time.Since(start)),
	)
	return catalog, nil
}
