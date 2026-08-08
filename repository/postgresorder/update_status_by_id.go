package postgresorder

import (
	"context"
	"telegram-service-platform/pkg/msgerror"

	"telegram-service-platform/entity"
	"telegram-service-platform/pkg/richerror"
)

func (d *DB) UpdateStatusByID(ctx context.Context, orderID uint64, status entity.Status) error {
	const Op = "postgresorder.UpdateStatusByID"

	query := `
	UPDATE orders
	SET
		status=$1,
		updated_at=NOW()
	WHERE id=$2
	`

	_, err := d.Pool.Connection().
		Exec(
			ctx,
			query,
			status,
			orderID,
		)

	if err != nil {
		return richerror.New(Op, err).WithKind(richerror.KindQueryFailure).WithMessage(msgerror.QueryFailed)
	}

	return nil
}
