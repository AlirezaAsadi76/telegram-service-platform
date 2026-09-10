package postgrespayment

import (
	"context"
	"errors"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"

	"github.com/jackc/pgx/v5"
)

func (d *DB) Fail(ctx context.Context, paymentID uint64) error {
	const Op = "postgrespayment.fail"

	const query = `
		UPDATE payments
		SET
			status = $1,
			updated_at = NOW()
		WHERE id = $2
		  AND status = $3
	`

	tag, err := d.executor.Exec(
		ctx,
		query,
		paymententity.PaymentStatusFailed,
		paymentID,
		paymententity.PaymentStatusPending,
	)
	if err != nil {
		return richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryFailed)
	}

	if tag.RowsAffected() == 1 {
		return nil
	}

	var status paymententity.PaymentStatus

	const selectQuery = `
		SELECT status
		FROM payments
		WHERE id = $1
	`

	err = d.executor.QueryRow(
		ctx,
		selectQuery,
		paymentID,
	).Scan(&status)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return richerror.New(Op, err).
				WithKind(richerror.KindNotFound).
				WithCode(richerror.CodePaymentNotFound)
		}

		return richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryFailed)
	}

	if status == paymententity.PaymentStatusFailed {
		return nil
	}

	return richerror.New(Op, nil).
		WithKind(richerror.KindConflict).
		WithCode(richerror.CodePaymentInvalidState)
}
