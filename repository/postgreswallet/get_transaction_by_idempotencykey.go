package postgreswallet

import (
	"context"
	"errors"
	"telegram-service-platform/entity/walletentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"

	"github.com/jackc/pgx/v5"
)

func (d *DB) GetTransactionByIdempotencyKey(ctx context.Context, key string) (*walletentity.WalletTransaction, error) {
	const Op = "postgreswallet.GetTransactionByIdempotencyKey"

	query := `
		SELECT id, wallet_id, user_id, type, amount, status, reference_id, idempotency_key, created_at, updated_at
		FROM wallet_transactions WHERE idempotency_key = $1
	`

	row := d.executor.QueryRow(ctx, query, key)

	walletTr, sErr := scanWalletTransaction(row)

	if sErr != nil {
		if errors.Is(sErr, pgx.ErrNoRows) {
			return nil, richerror.New(Op, sErr).
				WithKind(richerror.KindNotFound).
				WithMessage(msgerror.WalletTransactionNotFound)
		}

		return nil, richerror.New(Op, sErr).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryScanFailed)
	}

	return &walletTr, nil
}
