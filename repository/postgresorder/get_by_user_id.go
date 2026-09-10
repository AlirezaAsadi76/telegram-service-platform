package postgresorder

import (
	"context"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/logger"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"

	"go.uber.org/zap"
)

func (d *DB) GetByUserID(ctx context.Context, userID uint64) ([]*orderentity.Order, error) {
	const Op = "postgresorder.GetByUserID"

	query := `
	SELECT
		id,
		user_id,
		product_type,
		product_id,
		quantity,
		target_link,
		amount,
		currency,
		status,
		external_order_id,
		provider_id,
		metadata,
		created_at,
		updated_at
	FROM orders
	WHERE user_id=$1
	ORDER BY created_at DESC
	`

	rows, qErr := d.executor.Query(ctx, query, userID)

	if qErr != nil {
		logger.Logger.Debug("GetByUserID", zap.Error(qErr))
		return nil, richerror.New(Op, qErr).WithKind(richerror.KindQueryFailure).WithMessage(msgerror.QueryFailed)
	}

	defer rows.Close()

	orders := make([]*orderentity.Order, 0)

	for rows.Next() {

		order, sErr := scanOrder(rows)
		if sErr != nil {
			logger.Logger.Error("GetByUserID-scan", zap.Error(sErr))
			return nil, richerror.New(Op, sErr).WithKind(richerror.KindScanFailure).WithMessage(msgerror.QueryScanFailed)
		}

		orders = append(
			orders,
			&order,
		)
	}

	return orders, nil
}
