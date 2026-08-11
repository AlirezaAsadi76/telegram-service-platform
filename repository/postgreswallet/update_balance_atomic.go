package postgreswallet

import (
	"context"
	"errors"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (d *DB) UpdateBalanceAtomic(ctx context.Context, walletID uint64, newBalance, newVersion int64) error {
	const Op = "postgreswallet.UpdateBalanceAtomic"

	query := `
		UPDATE wallets 
		SET balance = $1, version = $2, updated_at = NOW()
		WHERE id = $3 AND version = $4
	`
	tag, err := d.Pool.Connection().Exec(ctx, query, newBalance, newVersion, walletID, newVersion-1)
	if err != nil {
		return richerror.New(Op, err).WithKind(richerror.KindQueryFailure).WithMessage(msgerror.QueryFailed)
	}
	if tag.RowsAffected() == 0 {
		return richerror.New(Op, errors.New("ErrConcurrentUpdate")).WithKind(richerror.KindQueryFailure).WithMessage(msgerror.QueryFailed)
	}
	return nil
}
