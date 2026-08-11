package postgreswallet

import (
	"context"
	"telegram-service-platform/entity/walletentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (d *DB) GetByUserID(ctx context.Context, userID uint64) (*walletentity.Wallet, error) {
	const Op = "postgreswallet.GetByUserID"

	query := `
			SELECT id, user_id, balance, currency, version, created_at, updated_at
			FROM wallets WHERE user_id = $1;
`
	row := d.Pool.Connection().QueryRow(ctx, query, userID)

	wallet, sErr := scanWallet(row)
	if sErr != nil {
		return nil, richerror.New(Op, sErr).WithKind(richerror.KindQueryFailure).WithMessage(msgerror.QueryScanFailed)
	}

	return &wallet, nil

}
