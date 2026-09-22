package postgresorder

import (
	"context"
	"time"

	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (d *DB) GetStaleProcessingWithoutEvidence(ctx context.Context, olderThan time.Duration, limit int) ([]*orderentity.Order, error) {
	const Op = "postgresorder.GetStaleProcessingWithoutEvidence"

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
		WHERE status = $1
		  AND updated_at <= NOW() - ($2 * INTERVAL '1 second')
		  AND COALESCE(external_order_id, '') = ''
		  AND NOT EXISTS (
				SELECT 1
				FROM order_fulfillment_attempts a
				WHERE a.order_id = orders.id
				  AND a.resolved_at IS NULL
		  )
		ORDER BY updated_at ASC
		LIMIT $3
	`

	rows, qErr := d.executor.Query(
		ctx,
		query,
		orderentity.OrderStatusProcessing,
		int64(olderThan.Seconds()),
		limit,
	)
	if qErr != nil {
		return nil, richerror.New(Op, qErr).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryFailed)
	}

	defer rows.Close()

	orders := make([]*orderentity.Order, 0, limit)

	for rows.Next() {
		order, scanErr := scanOrder(rows)
		if scanErr != nil {
			return nil, richerror.New(Op, scanErr).
				WithKind(richerror.KindScanFailure).
				WithMessage(msgerror.QueryScanFailed)
		}

		orders = append(orders, &order)
	}

	if err := rows.Err(); err != nil {
		return nil, richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryFailed)
	}

	return orders, nil
}
