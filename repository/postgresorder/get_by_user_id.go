package postgresorder

import (
	"context"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"

	"telegram-service-platform/entity"
)

func (d *DB) GetByUserID(ctx context.Context, userID uint64) ([]*entity.Order, error) {
	const Op = "postgresorder.GetByUserID"

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
	WHERE user_id=$1
	ORDER BY created_at DESC
	`

	rows, qErr := d.Pool.Connection().Query(ctx, query, userID)

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
