package postgrespayment

import (
	"context"
	"errors"
	"fmt"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"

	"github.com/jackc/pgx/v5"
)

func (d *DB) GetByExternalID(ctx context.Context, externalID string) (*paymententity.Payment, error) {
	const Op = "postgrespayment.getbyexternalid"

	const query = `
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
		WHERE external_id = $1
		LIMIT 1
	`

	row := d.executor.QueryRow(ctx, query, externalID)
	payment, err := scanPayment(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, richerror.New(Op, err).
				WithKind(richerror.KindNotFound).
				WithCode(richerror.CodePaymentNotFound)
		}
		fmt.Println(err.Error())
		return nil, richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryFailed)
	}

	return &payment, nil
}
