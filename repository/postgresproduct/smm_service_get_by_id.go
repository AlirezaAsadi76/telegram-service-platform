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

func (db *DB) SMMServiceGetByD(ctx context.Context, Id int64) (*smmentity.SMM, error) {
	const Op = "postgresproduct.SMMServiceGetByID"
	start := time.Now()

	query := `SELECT id, service_id, name, type, rate, min_quantity, max_quantity, drip_feed, refill, cancel, is_active, category, provider_name, created_at, updated_at
	          FROM smm_services WHERE id = $1`

	row := db.Pool.Connection().QueryRow(ctx, query, Id)
	smm, sErr := scanSMM(row)

	if sErr != nil {
		logger.Logger.Error("smm service get by service id failed",
			zap.String("op", Op),
			zap.Int64("service_id", Id),
			zap.Error(sErr),
			zap.Duration("duration", time.Since(start)),
		)
		return nil, richerror.New(Op, sErr).
			WithKind(richerror.KindNotFound).
			WithMessage(msgerror.ProductNotFound)
	}

	logger.Logger.Debug("smm service found by service id",
		zap.Int64("service_id", Id),
		zap.Duration("duration", time.Since(start)),
	)
	return &smm, nil
}
