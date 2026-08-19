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

// SMMMappingUpdate allows admin to edit an existing mapping (name, platform, category, etc.).
func (db *DB) SMMMappingUpdate(ctx context.Context, m *smmentity.SmmMapping) error {
	const Op = "postgresproduct.SMMMappingUpdate"
	start := time.Now()

	query := `
				UPDATE smm_service_mappings SET smm_service_id = $1, name = $2,
				                               platform = $3, category = $4,
				                               description = $5, is_active = $6,
				                               button_name = $7, sort_order = $8
				                               updated_at = NOW() WHERE id = $9`
	if _, err := db.Pool.Connection().Exec(ctx, query,
		m.SmmServiceId, m.Name, m.Platform, m.Category, m.Description, m.IsActive, m.ButtonName, m.SortOrder, m.Id,
	); err != nil {
		logger.Logger.Error("smm mapping update failed",
			zap.String("op", Op),
			zap.Int64("id", m.Id),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryFailed)
	}

	logger.Logger.Info("smm mapping updated",
		zap.Int64("id", m.Id),
		zap.String("name", m.Name),
		zap.Duration("duration", time.Since(start)),
	)
	return nil
}
