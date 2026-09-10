package postgrespayment

import (
	"context"
	"errors"
	"telegram-service-platform/entity/paymententity"

	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"

	"github.com/jackc/pgx/v5"
)

func (d *DB) Confirm(ctx context.Context, paymentID uint64) error {
	const Op = "postgrespay.confirm"

	return d.transactionProvider.Execute(ctx, func(tx pgx.Tx) error {
		var orderID uint64
		var status paymententity.PaymentStatus
		paymentSelectQuery := `SELECT status FROM payments WHERE id = $1`

		sErr := tx.QueryRow(ctx, paymentSelectQuery, paymentID).Scan(&status)
		if sErr != nil {
			if errors.Is(sErr, pgx.ErrNoRows) {
				return richerror.New(Op, sErr).
					WithKind(richerror.KindNotFound).
					WithCode(richerror.CodePaymentNotFound)
			}
			return richerror.New(Op, sErr).
				WithKind(richerror.KindQueryFailure).
				WithMessage(msgerror.QueryFailed)
		}
		if status == paymententity.PaymentStatusSuccess {
			return richerror.New(Op, sErr).
				WithKind(richerror.KindConflict).
				WithCode(richerror.CodePaymentAlreadyConfirmed)
		}

		paymentQuery := `
			UPDATE payments
			SET
				status = 'SUCCESS',
				updated_at = NOW()
			WHERE id = $1
			  AND status = 'PENDING'
			RETURNING order_id
		`

		err := tx.QueryRow(ctx, paymentQuery, paymentID).Scan(&orderID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return richerror.New(
					Op,
					err,
				).
					WithKind(richerror.KindConflict).
					WithMessage(msgerror.PaymentConfirmationConflict)
			}

			return richerror.New(Op, err).
				WithKind(richerror.KindQueryFailure).
				WithMessage(msgerror.QueryFailed)
		}

		orderQuery := `
			UPDATE orders
			SET
				status = 'PAID',
				updated_at = NOW()
			WHERE id = $1
			  AND status = 'PENDING'
		`

		tag, err := tx.Exec(ctx, orderQuery, orderID)
		if err != nil {
			return richerror.New(Op, err).
				WithKind(richerror.KindQueryFailure).
				WithMessage(msgerror.QueryFailed)
		}

		if tag.RowsAffected() != 1 {
			return richerror.New(
				Op,
				nil,
			).
				WithKind(richerror.KindConflict).
				WithMessage(msgerror.PaymentConfirmationConflict)
		}

		return nil
	})
}
