package paymentflowtesting

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func setPaymentExternalID(t *testing.T, pool *pgxpool.Pool, paymentID uint64, externalID string) {
	t.Helper()

	_, err := pool.Exec(
		context.Background(),
		`
			UPDATE payments
			SET external_id = $1
			WHERE id = $2
		`,
		externalID,
		paymentID,
	)

	if err != nil {
		t.Fatalf("set payment external ID: %v", err)
	}
}
