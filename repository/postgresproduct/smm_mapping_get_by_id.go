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

// SMMMappingGetByID returns a single mapping by its catalog ID.
// Use this during order creation to resolve the selected service.
func (db *DB) SMMMappingGetByID(ctx context.Context, id int64) (*smmentity.SmmMapping, error) {
	const Op = "postgresproduct.SMMMappingGetByID"
	start := time.Now()

	query := `SELECT id, smm_service_id, name, platform, category, description, is_active, created_at, updated_at
	          FROM smm_service_mappings WHERE id = $1`

	row := db.Pool.Connection().QueryRow(ctx, query, id)
	smm, sErr := scanSMMMapping(row)
	if sErr != nil {
		logger.Logger.Error("smm mapping get by id failed",
			zap.String("op", Op),
			zap.Int64("id", id),
			zap.Error(sErr),
			zap.Duration("duration", time.Since(start)),
		)
		return nil, richerror.New(Op, sErr).
			WithKind(richerror.KindNotFound).
			WithMessage(msgerror.ProductNotFound)
	}

	logger.Logger.Debug("smm mapping found by id",
		zap.Int64("id", id),
		zap.String("name", smm.Name),
		zap.Duration("duration", time.Since(start)),
	)
	return &smm, nil
}
