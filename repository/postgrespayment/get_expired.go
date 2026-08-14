package postgrespayment

import (
	"context"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (d *DB) GetExpired(ctx context.Context) ([]paymententity.Payment, error) {
	const Op = "postgrespayment.getExpired"

	query := `SELECT id, order_id, user_id, method, amount, currency, status, external_id, idempotency_key, callback_data, expired_at, created_at, updated_at
	          FROM payments WHERE status = $1 AND expired_at < NOW()`
	rows, qErr := d.Pool.Connection().Query(ctx, query, paymententity.PaymentStatusPending)
	if qErr != nil {
		return nil, richerror.New(Op, qErr).WithKind(richerror.KindQueryFailure).WithMessage(msgerror.QueryFailed)
	}
	defer rows.Close()

	var payments []paymententity.Payment
	for rows.Next() {
		payment, sErr := scanPayment(rows)
		if sErr != nil {
			return nil, richerror.New(Op, sErr).WithKind(richerror.KindScanFailure).WithMessage(msgerror.QueryScanFailed)
		}
		payments = append(payments, payment)
	}
	return payments, rows.Err()
}
