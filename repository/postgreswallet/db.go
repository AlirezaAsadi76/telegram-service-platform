package postgreswallet

import (
	"fmt"
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/walletentity"
	"telegram-service-platform/repository/postgres"

	"github.com/shopspring/decimal"
)

type DB struct {
	Pool *postgres.DB
}

func New(pool *postgres.DB) *DB {

	return &DB{Pool: pool}
}

func scanWallet(row postgres.Scanner) (walletentity.Wallet, error) {
	wallet := walletentity.Wallet{}
	var BalanceStr string
	err := row.Scan(
		&wallet.ID,
		&wallet.UserID,
		&BalanceStr,
		&wallet.Currency,
		&wallet.Version,
		&wallet.CreatedAt,
		&wallet.UpdatedAt)

	if err != nil {
		return wallet, err
	}

	Balance, nErr := decimal.NewFromString(BalanceStr)
	wallet.Balance = entity.Amount(Balance)
	if nErr != nil {
		return wallet, fmt.Errorf("failed to parse rate to decimal: %w", nErr)
	}
	return wallet, nil

}

func scanWalletTransaction(row postgres.Scanner) (walletentity.WalletTransaction, error) {
	walletTr := walletentity.WalletTransaction{}
	var AmountStr string

	err := row.Scan(
		&walletTr.ID,
		&walletTr.WalletID,
		&walletTr.UserID,
		&walletTr.Type,
		&AmountStr,
		&walletTr.Status,
		&walletTr.ReferenceID,
		&walletTr.IdempotencyKey,
		&walletTr.CreatedAt,
		&walletTr.UpdatedAt,
	)
	if err != nil {
		return walletTr, err
	}
	Amount, nErr := decimal.NewFromString(AmountStr)
	walletTr.Amount = entity.Amount(Amount)
	if nErr != nil {
		return walletTr, fmt.Errorf("failed to parse rate to decimal: %w", nErr)
	}

	return walletTr, err

}
