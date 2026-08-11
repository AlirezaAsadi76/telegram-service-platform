package postgreswallet

import (
	"context"
	"telegram-service-platform/entity/walletentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (d *DB) UpdateTransactionStatus(ctx context.Context, txID uint64, status walletentity.WalletTransactionStatus) error {
	const Op = "postgreswallet.UpdateTransactionStatus"
	query := `UPDATE wallet_transactions SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := d.Pool.Connection().Exec(ctx, query, status, txID)

	if err != nil {
		return richerror.New(Op, err).WithKind(richerror.KindQueryFailure).WithMessage(msgerror.QueryFailed)
	}
	return err
}
