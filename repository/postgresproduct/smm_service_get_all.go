package postgresproduct

import (
	"context"
	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/logger"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
	"time"

	"go.uber.org/zap"
)

// SMMServiceGetAll returns all raw services from JAP ordered by category and name.
func (db *DB) SMMServiceGetAll(ctx context.Context) ([]smmentity.SMM, error) {
	const Op = "postgresproduct.SMMServiceGetAll"
	start := time.Now()

	query := `SELECT id, service_id, name, type, rate, min_quantity, max_quantity, drip_feed, refill, cancel, is_active, category, provider_name, created_at, updated_at
	          FROM smm_services ORDER BY category, name`
	rows, err := db.Pool.Connection().Query(ctx, query)
	if err != nil {
		logger.Logger.Error("smm service get all failed",
			zap.String("op", Op),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return nil, richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryFailed)
	}
	defer rows.Close()

	var services []smmentity.SMM
	for rows.Next() {
		var s smmentity.SMM
		if err := rows.Scan(&s.Id, &s.Service, &s.Name, &s.Type, &s.Rate, &s.Min, &s.Max,
			&s.DripFeed, &s.Refill, &s.Cancel, &s.IsActive, &s.Category, &s.ProviderName, &s.CreatedAt, &s.UpdatedAt); err != nil {
			logger.Logger.Error("smm service scan failed",
				zap.String("op", Op),
				zap.Error(err),
			)
			return nil, richerror.New(Op, err).
				WithKind(richerror.KindQueryFailure).
				WithMessage(msgerror.QueryScanFailed)
		}
		services = append(services, s)
	}

	if err := rows.Err(); err != nil {
		logger.Logger.Error("smm service rows error",
			zap.String("op", Op),
			zap.Error(err),
		)
		return nil, richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryFailed)
	}

	logger.Logger.Debug("smm service get all completed",
		zap.Int("count", len(services)),
		zap.Duration("duration", time.Since(start)),
	)
	return services, nil
}
