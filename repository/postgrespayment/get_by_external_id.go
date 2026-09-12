package postgrespayment

import (
	"context"
	"errors"

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

	var payment paymententity.Payment

	err := d.executor.QueryRow(ctx, query, externalID).
		Scan(
			&payment.ID,
			&payment.OrderID,
			&payment.UserID,
			&payment.Method,
			&payment.Amount,
			&payment.Currency,
			&payment.Status,
			&payment.ExternalID,
			&payment.ProviderReferenceID,
			&payment.PaymentURL,
			&payment.IdempotencyKey,
			&payment.CallbackData,
			&payment.ExpiredAt,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, richerror.New(Op, err).
				WithKind(richerror.KindNotFound).
				WithCode(richerror.CodePaymentNotFound)
		}

		return nil, richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryFailed)
	}

	return &payment, nil
}
