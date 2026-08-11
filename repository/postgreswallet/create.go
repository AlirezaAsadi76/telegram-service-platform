package postgreswallet

import (
	"context"
	"telegram-service-platform/entity/walletentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (d *DB) Create(ctx context.Context, wallet *walletentity.Wallet) error {
	const Op = "postgreswallet.Create"
	query := `

		INSERT INTO wallets (user_id, balance, currency, version)
		VALUES ($1, $2, $3, 1)
		RETURNING id, created_at, updated_at
	`

	sErr := d.Pool.Connection().QueryRow(ctx, query,
		wallet.UserID, wallet.Balance, wallet.Currency,
	).Scan(&wallet.ID, &wallet.CreatedAt, &wallet.UpdatedAt)

	if sErr != nil {
		return richerror.New(Op, sErr).WithKind(richerror.KindQueryFailure).WithMessage(msgerror.QueryScanFailed)
	}

	return nil

}
