package postgresorder

import (
	"context"
	"telegram-service-platform/entity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (d *DB) GetByStatus(ctx context.Context, status entity.Status) ([]*entity.Order, error) {
	const Op = "postgresorder.GetByStatus"

	query := `
	SELECT
		id,
		user_id,
		product_type,
		product_id,
		quantity,
		amount,
		currency,
		status,
		metadata,
		created_at,
		updated_at
	FROM orders
	WHERE status=$1
	ORDER BY created_at DESC
	`

	rows, qErr := d.Pool.Connection().Query(ctx, query, status)

	if qErr != nil {
		return nil, richerror.New(Op, qErr).WithKind(richerror.KindQueryFailure).WithMessage(msgerror.QueryFailed)
	}

	defer rows.Close()

	orders := make([]*entity.Order, 0)

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
