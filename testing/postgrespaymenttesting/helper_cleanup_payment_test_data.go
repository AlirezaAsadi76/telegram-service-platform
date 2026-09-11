package postgrespaymenttesting

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func cleanupPaymentTestData(t *testing.T, pool *pgxpool.Pool, paymentID uint64, orderID uint64) {
	t.Helper()

	_, err := pool.Exec(
		context.Background(),
		`DELETE FROM payments WHERE id = $1`,
		paymentID,
	)
	if err != nil {
		t.Fatalf("delete test payment: %v", err)
	}

	_, err = pool.Exec(
		context.Background(),
		`DELETE FROM orders WHERE id = $1`,
		orderID,
	)
	if err != nil {
		t.Fatalf("delete test order: %v", err)
	}
}
