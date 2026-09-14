package paymentflowtesting

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func getPaymentIDByIdempotencyKey(t *testing.T, pool *pgxpool.Pool, idempotencyKey string) uint64 {
	t.Helper()

	var paymentID uint64

	err := pool.QueryRow(
		context.Background(),
		`
			SELECT id
			FROM payments
			WHERE idempotency_key = $1
		`,
		idempotencyKey,
	).Scan(&paymentID)
	if err != nil {
		t.Fatalf("get payment id by idempotency key: %v", err)
	}

	return paymentID
}
