package postgresproduct

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

func (db *DB) SMMServiceCreateOrUpdate(ctx context.Context, s smmentity.SMM) error {
	const Op = "postgresproduct.SMMServiceCreateOrUpdate"
	start := time.Now()

	query := `
		INSERT INTO smm_services (service_id, name, type, rate, min_quantity, max_quantity, drip_feed, refill, cancel, is_active, category, provider_name)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (service_id) DO UPDATE SET
			name = EXCLUDED.name,
			type = EXCLUDED.type,
			rate = EXCLUDED.rate,
			min_quantity = EXCLUDED.min_quantity,
			max_quantity = EXCLUDED.max_quantity,
			drip_feed = EXCLUDED.drip_feed,
			refill = EXCLUDED.refill,
			cancel = EXCLUDED.cancel,
			is_active = EXCLUDED.is_active,
			category = EXCLUDED.category,
			provider_name = EXCLUDED.provider_name,
			updated_at = NOW()
	`
	_, err := db.Pool.Connection().Exec(ctx, query,
		s.Service, s.Name, s.Type, s.Rate, s.Min, s.Max,
		s.DripFeed, s.Refill, s.Cancel, s.IsActive, s.Category, s.ProviderName,
	)

	if err != nil {
		metrics.SMMProviderRequests.WithLabelValues("db_upsert", "error").Inc()
		logger.Logger.Error("smm service upsert failed",
			zap.String("op", Op),
			zap.Int64("service_id", s.Service),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryFailed)
	}

	metrics.SMMProviderRequests.WithLabelValues("db_upsert", "success").Inc()
	logger.Logger.Debug("smm service upserted",
		zap.Int64("service_id", s.Service),
		zap.String("name", s.Name),
		zap.Duration("duration", time.Since(start)),
	)
	return nil
}
