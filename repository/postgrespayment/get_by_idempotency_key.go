package postgrespayment

import (
	"context"
	"errors"

	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"

	"github.com/jackc/pgx/v5"
)

func (d *DB) GetByIdempotencyKey(ctx context.Context, key string) (*paymententity.Payment, error) {
	const op = "postgrespayment.getbyidempotencykey"

	query := `
		SELECT
			id,
			order_id,
			user_id,
			method,
			amount,
			currency,
			status,
			external_id,
			provider_reference_id,
			payment_url,
			idempotency_key,
			callback_data,
			expired_at,
			created_at,
			updated_at
		FROM payments
		WHERE idempotency_key = $1
		LIMIT 1
	`

	row := d.executor.QueryRow(ctx, query, key)

	payment, err := scanPayment(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, richerror.New(op, err).
				WithKind(richerror.KindNotFound).
				WithCode(richerror.CodePaymentNotFound)
		}

		return nil, richerror.New(op, err).
			WithKind(richerror.KindScanFailure).
			WithMessage(msgerror.QueryScanFailed)
	}

	return &payment, nil
}
