package postgresorder

import (
	"context"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (d *DB) UpdateStatus(ctx context.Context, id uint64, status orderentity.OrderStatus, externalOrderID string, providerID *uint64) error {
	const Op = "postgresorder.UpdateStatus"

	query := `
		UPDATE orders 
		SET status = $1, external_order_id = $2, provider_id = $3, updated_at = NOW()
		WHERE id = $4
	`
	_, err := d.Pool.Connection().Exec(ctx, query, status, externalOrderID, providerID, id)

	if err != nil {
		return richerror.New(Op, err).WithKind(richerror.KindQueryFailure).WithMessage(msgerror.QueryFailed)
	}

	return err
}
