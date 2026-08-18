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

// SMMMappingGetByPlatformCategory returns active mappings for a specific platform and category.
// Use this when user has already selected a platform and category in the bot.
func (db *DB) SMMMappingGetByPlatformCategory(ctx context.Context, platform smmentity.PlatformType, category smmentity.Category) ([]smmentity.SmmMapping, error) {
	const Op = "postgresproduct.SMMMappingGetByPlatformCategory"
	start := time.Now()

	query := `SELECT id, smm_service_id, name, platform, category, description, is_active, button_name, created_at, updated_at
	          FROM smm_service_mappings m WHERE platform = $1 AND category = $2 AND is_active = true ORDER BY name`
	rows, gErr := db.Pool.Connection().Query(ctx, query, platform, category)
	if gErr != nil {
		logger.Logger.Error("smm mapping get by platform category failed",
			zap.String("op", Op),
			zap.String("platform", string(platform)),
			zap.String("category", string(category)),
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

	logger.Logger.Debug("smm mapping get by platform category completed",
		zap.String("platform", string(platform)),
		zap.String("category", string(category)),
		zap.Int("count", len(mappings)),
		zap.Duration("duration", time.Since(start)),
	)
	return mappings, nil
}
