package postgresorder

import (
	"context"

	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (d *DB) MarkFulfillmentAttemptResolved(ctx context.Context, attemptID uint64) error {
	const Op = "postgresorder.MarkFulfillmentAttemptResolved"

	query := `
		UPDATE order_fulfillment_attempts
		SET resolved_at = NOW(),
			updated_at = NOW()
		WHERE id = $1
		  AND resolved_at IS NULL
	`

	_, err := d.executor.Exec(ctx, query, attemptID)
	if err != nil {
		return richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryFailed)
	}

	return nil
}
