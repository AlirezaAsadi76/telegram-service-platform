package postgresorder

import (
	"context"

	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (d *DB) GetByIDForUpdate(ctx context.Context, orderID uint64) (*orderentity.Order, error) {
	const Op = "postgresorder.GetByIDForUpdate"

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
		WHERE id = $1
		FOR UPDATE
	`

	row := d.executor.QueryRow(ctx, query, orderID)

	order, err := scanOrder(row)
	if err != nil {
		return nil, richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryScanFailed)
	}

	return &order, nil
}
