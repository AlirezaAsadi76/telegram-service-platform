package postgreswallet

import (
	"context"
	"telegram-service-platform/entity/walletentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (d *DB) GetTransactionByIdempotencyKey(ctx context.Context, key string) (*walletentity.WalletTransaction, error) {
	const Op = "postgreswallet.GetTransactionByIdempotencyKey"

	query := `
		SELECT id, wallet_id, user_id, type, amount, status, reference_id, idempotency_key, created_at
		FROM wallet_transactions WHERE idempotency_key = $1
	`

	row := d.Pool.Connection().QueryRow(ctx, query, key)

	walletTr, sErr := scanWalletTransaction(row)

	if sErr != nil {
		return nil, richerror.New(Op, sErr).WithKind(richerror.KindQueryFailure).WithMessage(msgerror.QueryScanFailed)
	}

	return &walletTr, nil
}
