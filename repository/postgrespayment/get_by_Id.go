package postgrespayment

import (
	"context"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (d *DB) GetByID(ctx context.Context, id uint64) (*paymententity.Payment, error) {
	const Op = "postgrespayment.getbyid"

	query := `
		SELECT id, order_id, user_id, method, amount, currency, status, external_id, idempotency_key, callback_data, expired_at, created_at, updated_at
		FROM payments WHERE id = $1
	`

	row := d.Pool.Connection().QueryRow(ctx, query, id)

	payment, psErr := scanPayment(row)
	if psErr != nil {
		return nil, richerror.New(Op, psErr).WithKind(richerror.KindScanFailure).WithMessage(msgerror.QueryScanFailed)
	}

	return &payment, nil
}
