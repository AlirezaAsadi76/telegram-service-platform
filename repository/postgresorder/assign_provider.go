package postgresorder

import (
	"context"

	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (d *DB) AssignProvider(ctx context.Context, orderID uint64, providerID uint64) error {
	const Op = "postgresorder.AssignProvider"

	query := `
		UPDATE orders
		SET provider_id = $1,
		    updated_at = NOW()
		WHERE id = $2
		  AND status = 'PROCESSING'
		  AND (
			  provider_id IS NULL
			  OR provider_id = $1
		  )
	`

	tag, err := d.executor.Exec(ctx, query, providerID, orderID)
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
