package postgrescheckouttesting

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

func getWalletBalance(t *testing.T, pool *pgxpool.Pool, walletID uint64) decimal.Decimal {
	t.Helper()

	var balanceText string

	err := pool.QueryRow(
		context.Background(),
		`SELECT balance::text FROM wallets WHERE id = $1`,
		walletID,
	).Scan(&balanceText)

	if err != nil {
		t.Fatalf(
			"get wallet balance: %v",
			err,
		)
	}

	balance, err := decimal.NewFromString(balanceText)
	if err != nil {
		t.Fatalf(
			"parse wallet balance %q: %v",
			balanceText,
			err,
		)
	}

	return balance
}
