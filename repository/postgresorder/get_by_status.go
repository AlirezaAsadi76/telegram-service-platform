package postgresorder

import (
	"context"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (d *DB) GetByStatus(ctx context.Context, status orderentity.OrderStatus) ([]*orderentity.Order, error) {
	const Op = "postgresorder.GetByStatus"

	query := `
		SELECT id, user_id, product_type, product_id, quantity, target_link, amount, currency, status,
		       external_order_id, provider_id, metadata, created_at, updated_at
		FROM orders WHERE status = $1
	`

	rows, qErr := d.Pool.Connection().Query(ctx, query, status)
	if qErr != nil {
		return nil, richerror.New(Op, qErr).WithKind(richerror.KindQueryFailure).WithMessage(msgerror.QueryFailed)
	}

	defer rows.Close()

	orders := make([]*orderentity.Order, 0)

	for rows.Next() {

		order, sErr := scanOrder(rows)
		if sErr != nil {
			return nil, richerror.New(Op, sErr).WithKind(richerror.KindScanFailure).WithMessage(msgerror.QueryScanFailed)
		}

		orders = append(
			orders,
			&order,
		)
	}

	return orders, nil
}
