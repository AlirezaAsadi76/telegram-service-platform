package paymentflowtesting

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func countPaymentByIdempotencyKey(t *testing.T, pool *pgxpool.Pool, idempotencyKey string) int64 {
	t.Helper()

	var count int64
	err := pool.QueryRow(context.Background(),
		"SELECT count(*) FROM payments WHERE idempotency_key = $1", idempotencyKey).Scan(&count)
	if err != nil {
		t.Fatalf("check exist payment: %v", err)
	}
	return count
}
