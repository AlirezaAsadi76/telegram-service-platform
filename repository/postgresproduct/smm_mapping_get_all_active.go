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

// SMMMappingGetAllActive returns all active mappings ordered by platform, category, name.
// Use this to build the full catalog menu in Telegram bot.
func (db *DB) SMMMappingGetAllActive(ctx context.Context) ([]smmentity.SmmMapping, error) {
	const Op = "postgresproduct.SMMMappingGetAllActive"
	start := time.Now()

	query := `SELECT m.id, m.smm_service_id, m.name, m.platform, m.category, m.description, m.is_active, m.created_at, m.updated_at
	          FROM smm_service_mapping m WHERE m.is_active = true ORDER BY m.platform, m.category, m.name`
	rows, gErr := db.Pool.Connection().Query(ctx, query)
	if gErr != nil {
		logger.Logger.Error("smm mapping get all active failed",
			zap.String("op", Op),
			zap.Error(gErr),
			zap.Duration("duration", time.Since(start)),
		)
		return nil, richerror.New(Op, gErr).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryFailed)
	}

	defer rows.Close()
	var mappings []smmentity.SmmMapping
	for rows.Next() {

		smm, sErr := scanSMMMapping(rows)
		if sErr != nil {
			logger.Logger.Error("smm mapping scan failed",
				zap.String("op", Op),
				zap.Error(sErr),
			)
			return nil, richerror.New(Op, sErr).
				WithKind(richerror.KindScanFailure).
				WithMessage(msgerror.QueryScanFailed)
		}
		mappings = append(mappings, smm)
	}

	if err := rows.Err(); err != nil {
		logger.Logger.Error("smm mapping rows error",
			zap.String("op", Op),
			zap.Error(err),
		)
		return nil, richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryFailed)
	}

	logger.Logger.Debug("smm mapping get all active completed",
		zap.Int("count", len(mappings)),
		zap.Duration("duration", time.Since(start)),
	)
	return mappings, nil
}

func scanSMMMapping(row postgres.Scanner) (smmentity.SmmMapping, error) {
	var m smmentity.SmmMapping
	err := row.Scan(
		&m.Id,
		&m.SmmServiceId,
		&m.Name,
		&m.Platform,
		&m.Category,
		&m.Description,
		&m.IsActive,
		&m.ButtonName,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	return m, err
}
