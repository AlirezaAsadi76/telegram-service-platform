package postgresorder

import (
	"context"

	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (d *DB) GetUnresolvedFulfillmentAttempts(ctx context.Context, limit int) ([]orderentity.FulfillmentAttempt, error) {
	const Op = "postgresorder.GetUnresolvedFulfillmentAttempts"

	query := `
		SELECT
			id,
			order_id,
			provider_id,
			outcome,
			external_order_id,
			resolved_at,
			created_at,
			updated_at
		FROM order_fulfillment_attempts
		WHERE resolved_at IS NULL
		ORDER BY id ASC
		LIMIT $1
	`

	rows, err := d.executor.Query(ctx, query, limit)
	if err != nil {
		return nil, richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryFailed)
	}
	defer rows.Close()

	attempts := make([]orderentity.FulfillmentAttempt, 0, limit)

	for rows.Next() {
		attempt, sErr := scanFulfillmentAttempt(rows)
		if sErr != nil {
			return nil, richerror.New(Op, sErr).
				WithKind(richerror.KindQueryFailure).
				WithMessage(msgerror.QueryScanFailed)
		}
		attempts = append(attempts, attempt)
	}

	if err := rows.Err(); err != nil {
		return nil, richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryFailed)
	}

	return attempts, nil
}
