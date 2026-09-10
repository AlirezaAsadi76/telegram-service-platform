package postgrespayment

import (
	"context"

	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (d *DB) MarkInitiated(
	ctx context.Context,
	paymentID uint64,
	status paymententity.PaymentStatus,
	externalID string,
	paymentURL string,
) error {
	const op = "postgrespayment.mark_initiated"

	query := `
		UPDATE payments
		SET
			status = $1,
			external_id = $2,
			payment_url = $3,
			updated_at = NOW()
		WHERE id = $4
		  AND status = 'CREATING'
	`

	result, err := d.executor.Exec(
		ctx,
		query,
		status,
		externalID,
		paymentURL,
		paymentID,
	)
	if err != nil {
		return richerror.New(op, err).
			WithKind(richerror.KindQueryFailure).
			WithCode(richerror.CodePaymentInitiationUpdateFailed).
			WithMessage(msgerror.QueryFailed)
	}

	if result.RowsAffected() == 0 {
		return richerror.New(op, nil).
			WithKind(richerror.KindConflict).
			WithCode(richerror.CodePaymentInvalidState)
	}

	return nil
}
