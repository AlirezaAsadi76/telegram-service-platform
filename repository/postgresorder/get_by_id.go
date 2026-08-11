package postgresorder

import (
	"context"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/pkg/msgerror"

	"telegram-service-platform/pkg/richerror"
)

func (d *DB) GetByID(ctx context.Context, orderID uint64) (*orderentity.Order, error) {
	const Op = "postgresorder.GetByID"

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
	WHERE id=$1
	`

	row := d.Pool.Connection().QueryRow(ctx, query, orderID)

	order, oErr := scanOrder(row)
	if oErr != nil {
		return nil, richerror.New(Op, oErr).WithKind(richerror.KindScanFailure).WithMessage(msgerror.QueryScanFailed)
	}

	return &order, nil
}
