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

// SMMMappingCreate inserts a new curated mapping for the bot catalog.
func (db *DB) SMMMappingCreate(ctx context.Context, m *smmentity.SmmMapping) error {
	const Op = "postgresproduct.SMMMappingCreate"
	start := time.Now()

	query := `INSERT INTO smm_service_mapping (smm_service_id, name, platform, category, description, is_active, button_name)
	          VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at`
	if err := db.Pool.Connection().QueryRow(ctx, query,
		m.SmmServiceId, m.Name, m.Platform, m.Category, m.Description, m.IsActive, m.ButtonName,
	).Scan(&m.Id, &m.CreatedAt, &m.UpdatedAt); err != nil {
		logger.Logger.Error("smm mapping create failed",
			zap.String("op", Op),
			zap.String("name", m.Name),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryFailed)
	}

	logger.Logger.Info("smm mapping created",
		zap.Int64("id", m.Id),
		zap.String("name", m.Name),
		zap.String("platform", m.Platform),
		zap.String("category", m.Category),
		zap.Duration("duration", time.Since(start)),
	)
	return nil
}
