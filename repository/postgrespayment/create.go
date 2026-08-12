package postgrespayment

import (
	"context"
	"encoding/json"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"

	"telegram-service-platform/entity/paymententity"
)

func (d *DB) Create(ctx context.Context, payment *paymententity.Payment) error {
	const Op = "postgrespayment.Create"

	metadata, err := json.Marshal(payment.CallbackData)
	if err != nil {
		return richerror.New(Op, err).WithKind(richerror.KindQueryFailure).WithMessage(msgerror.QueryFailed)
	}

	query := `
		INSERT INTO payments (order_id, user_id, method, amount, currency, status, external_id, idempotency_key, callback_data, expired_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`
	qErr := d.Pool.Connection().QueryRow(ctx, query,
		payment.OrderID, payment.UserID, payment.Method, payment.Amount,
		payment.Currency, payment.Status, payment.ExternalID,
		payment.IdempotencyKey, metadata, payment.ExpiredAt,
	).Scan(&payment.ID, &payment.CreatedAt, &payment.UpdatedAt)

	if qErr != nil {
		return richerror.New(Op, err).WithKind(richerror.KindQueryFailure).WithMessage(msgerror.QueryFailed)
	}

	return nil
}
