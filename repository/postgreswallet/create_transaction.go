package postgreswallet

import (
	"context"
	"telegram-service-platform/entity/walletentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (d *DB) CreateTransaction(ctx context.Context, tx *walletentity.WalletTransaction) error {
	const Op = "postgreswallet.CreateTransaction"

	query := `
		INSERT INTO wallet_transactions 
		(wallet_id, user_id, type, amount, status, reference_id, idempotency_key)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	sErr := d.Pool.Connection().QueryRow(ctx, query,
		tx.WalletID, tx.UserID, tx.Type, tx.Amount, tx.Status,
		tx.ReferenceID, tx.IdempotencyKey,
	).Scan(&tx.ID, &tx.CreatedAt, &tx.UpdatedAt)

	if sErr != nil {
		return richerror.New(Op, sErr).WithKind(richerror.KindQueryFailure).WithMessage(msgerror.QueryScanFailed)
	}

	return nil
}
