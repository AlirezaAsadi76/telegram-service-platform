package postgresorder

import (
	"context"

	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (d *DB) CompleteProcessing(ctx context.Context, orderID uint64) (bool, error) {
	const Op = "postgresorder.CompleteProcessing"

	query := `
		UPDATE orders
		SET status = 'COMPLETED',
		    updated_at = NOW()
		WHERE id = $1
		  AND status = 'PROCESSING'
	`

	tag, err := d.executor.Exec(ctx, query, orderID)
	if err != nil {
		return false, richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryFailed)
	}

	return tag.RowsAffected() == 1, nil
}
