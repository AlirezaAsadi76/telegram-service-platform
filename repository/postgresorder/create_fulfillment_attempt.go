package postgresorder

import (
	"context"
	"errors"

	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"

	"github.com/jackc/pgx/v5"
)

func (d *DB) CreateFulfillmentAttempt(ctx context.Context, attempt *orderentity.FulfillmentAttempt) (uint64, error) {
	const Op = "postgresorder.CreateFulfillmentAttempt"

	query := `
		INSERT INTO order_fulfillment_attempts (
			order_id,
			provider_id,
			outcome,
			external_order_id,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		ON CONFLICT (order_id) DO NOTHING
		RETURNING id
	`

	var id uint64
	if err := d.executor.QueryRow(
		ctx,
		query,
		attempt.OrderID,
		attempt.ProviderID,
		attempt.Outcome,
		attempt.ExternalOrderID,
	).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, richerror.New(Op, nil).
				WithKind(richerror.KindConflict).
				WithMessage(msgerror.OrderUpdateFailed)
		}

		return 0, richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryFailed)
	}

	return id, nil
}
