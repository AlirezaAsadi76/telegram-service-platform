package postgresproduct

import (
	"context"
	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/logger"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
	"telegram-service-platform/repository/postgres"
	"time"

	"go.uber.org/zap"
)

// SMMServiceGetAll returns all raw services from JAP ordered by category and name.
func (db *DB) SMMServiceGetAll(ctx context.Context) ([]smmentity.SMM, error) {
	const Op = "postgresproduct.SMMServiceGetAll"
	start := time.Now()

	query := `SELECT id, service_id, name, type, rate, min_quantity, max_quantity, drip_feed, refill, cancel, is_active, category, provider_name, created_at, updated_at
	          FROM smm_services ORDER BY category, name`
	rows, gErr := db.Pool.Connection().Query(ctx, query)
	if gErr != nil {
		logger.Logger.Error("smm service get all failed",
			zap.String("op", Op),
			zap.Error(gErr),
			zap.Duration("duration", time.Since(start)),
		)
		return nil, richerror.New(Op, gErr).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryFailed)
	}
	defer rows.Close()

	var services []smmentity.SMM
	for rows.Next() {

		smm, sErr := scanSMM(rows)
		if sErr != nil {
			logger.Logger.Error("smm service scan failed",
				zap.String("op", Op),
				zap.Error(sErr),
			)
			return nil, richerror.New(Op, sErr).
				WithKind(richerror.KindScanFailure).
				WithMessage(msgerror.QueryScanFailed)
		}

		services = append(services, smm)
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

func scanSMM(row postgres.Scanner) (smmentity.SMM, error) {
	var s smmentity.SMM
	err := row.Scan(
		&s.Id,
		&s.Service,
		&s.Name,
		&s.Type,
		&s.Rate,
		&s.Min,
		&s.Max,
		&s.DripFeed,
		&s.Refill,
		&s.Cancel,
		&s.IsActive,
		&s.Category,
		&s.ProviderName,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	return s, err
}
