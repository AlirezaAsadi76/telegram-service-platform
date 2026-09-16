package postgresorder

import (
	"context"

	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (d *DB) SetProviderOrder(ctx context.Context, orderID uint64, providerID uint64, externalOrderID string) error {
	const Op = "postgresorder.SetProviderOrder"

	query := `
		UPDATE orders
		SET provider_id = $1,
		    external_order_id = $2,
		    updated_at = NOW()
		WHERE id = $3
		  AND status = 'PROCESSING'
	`

	tag, err := d.executor.Exec(
		ctx,
		query,
		providerID,
		externalOrderID,
		orderID,
	)
	if err != nil {
		return richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryFailed)
	}

	if tag.RowsAffected() != 1 {
		return richerror.New(Op, nil).
			WithKind(richerror.KindConflict).
			WithMessage(msgerror.OrderUpdateFailed)
	}

	return nil
}
