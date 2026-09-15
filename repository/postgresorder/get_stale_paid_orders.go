package postgresorder

import (
	"context"
	"time"

	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (d *DB) GetStalePaid(ctx context.Context, olderThan time.Duration, limit int) ([]*orderentity.Order, error) {
	const Op = "postgresorder.GetStalePaid"

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
		ORDER BY updated_at ASC
		LIMIT $3
	`

	rows, qErr := d.executor.Query(
		ctx,
		query,
		orderentity.OrderStatusPaid,
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
		order, err := scanOrder(rows)
		if err != nil {
			return nil, richerror.New(Op, err).
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
