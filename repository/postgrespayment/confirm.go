package postgrespayment

import (
	"context"
	"errors"
	"telegram-service-platform/entity/paymententity"

	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"

	"github.com/jackc/pgx/v5"
)

func (d *DB) Confirm(ctx context.Context, paymentID uint64, providerReferenceID string) error {
	const Op = "postgrespay.confirm"

	return d.transactionProvider.Execute(ctx, func(tx pgx.Tx) error {
		var orderID uint64
		var status paymententity.PaymentStatus

		const paymentSelectQuery = `SELECT order_id, status FROM payments WHERE id = $1 FOR UPDATE`

		sErr := tx.QueryRow(ctx, paymentSelectQuery, paymentID).Scan(&orderID, &status)
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
			return richerror.New(Op, nil).
				WithKind(richerror.KindConflict).
				WithCode(richerror.CodePaymentAlreadyConfirmed)
		}

		if status != paymententity.PaymentStatusPending {
			return richerror.New(Op, nil).
				WithKind(richerror.KindConflict).
				WithCode(richerror.CodePaymentInvalidState)
		}

		const paymentUpdateQuery = `
			UPDATE payments
			SET
				status = 'SUCCESS',
				provider_reference_id = $1,
				updated_at = NOW()
			WHERE id = $2
			  AND status = 'PENDING'
		`

		tag, err := tx.Exec(ctx, paymentUpdateQuery, providerReferenceID, paymentID)
		if err != nil {
			return richerror.New(Op, err).
				WithKind(richerror.KindQueryFailure).
				WithMessage(msgerror.QueryFailed)
		}

		if tag.RowsAffected() != 1 {
			return richerror.New(Op, nil).
				WithKind(richerror.KindConflict).
				WithCode(richerror.CodePaymentInvalidState)
		}

		const orderUpdateQuery = `
			UPDATE orders
			SET
				status = 'PAID',
				updated_at = NOW()
			WHERE id = $1
			  AND status = 'PENDING'
		`

		tag, err = tx.Exec(ctx, orderUpdateQuery, orderID)
		if err != nil {
			return richerror.New(Op, err).
				WithKind(richerror.KindQueryFailure).
				WithMessage(msgerror.QueryFailed)
		}

		if tag.RowsAffected() != 1 {
			return richerror.New(Op, nil).
				WithKind(richerror.KindConflict).
				WithCode(richerror.CodeOrderInvalidState)
		}

		return nil
	})
}
