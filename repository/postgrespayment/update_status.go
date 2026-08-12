package postgrespayment

import (
	"context"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (d *DB) UpdateStatus(ctx context.Context, id uint64, status paymententity.PaymentStatus) error {
	const Op = "postgrespay.updatestatus"
	query := `UPDATE payments SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := d.Pool.Connection().Exec(ctx, query, status, id)

	if err != nil {
		return richerror.New(Op, err).WithKind(richerror.KindQueryFailure).WithMessage(msgerror.QueryFailed)
	}

	return err
}
