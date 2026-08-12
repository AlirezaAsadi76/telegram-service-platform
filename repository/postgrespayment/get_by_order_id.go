package postgrespayment

import (
	"context"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (d *DB) GetByOrderID(ctx context.Context, orderID uint64) (*paymententity.Payment, error) {
	const Op = "postgrespay.getbyorderid"
	query := `
		SELECT id, order_id, user_id, method, amount, currency, status, external_id, idempotency_key,CallbackData , ExpiredAt, created_at, updated_at
		FROM payments WHERE order_id = $1 ORDER BY created_at DESC LIMIT 1
	`

	row := d.Pool.Connection().QueryRow(ctx, query, orderID)

	payment, psErr := scanPayment(row)
	if psErr != nil {
		return nil, richerror.New(Op, psErr).WithKind(richerror.KindScanFailure).WithMessage(msgerror.QueryScanFailed)
	}

	return &payment, nil
}
