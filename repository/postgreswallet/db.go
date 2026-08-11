package postgreswallet

import (
	"telegram-service-platform/entity/walletentity"
	"telegram-service-platform/repository/postgres"
)

type DB struct {
	Pool *postgres.DB
}

func New(pool *postgres.DB) *DB {

	return &DB{Pool: pool}
}

func scanWallet(row postgres.Scanner) (walletentity.Wallet, error) {
	wallet := walletentity.Wallet{}

	err := row.Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.Balance,
		&wallet.Currency,
		&wallet.CreatedAt,
		&wallet.UpdatedAt)

	return wallet, err

}

func scanWalletTransaction(row postgres.Scanner) (walletentity.WalletTransaction, error) {
	walletTr := walletentity.WalletTransaction{}

	err := row.Scan(
		&walletTr.ID,
		&walletTr.WalletID,
		&walletTr.UserID,
		&walletTr.Amount,
		&walletTr.Status,
		&walletTr.ReferenceID,
		&walletTr.IdempotencyKey,
		&walletTr.CreatedAt,
	)

	return walletTr, err

}
