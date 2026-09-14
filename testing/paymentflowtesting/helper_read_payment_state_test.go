package paymentflowtesting

import (
	"context"
	"telegram-service-platform/entity/paymententity"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func readPaymentState(t *testing.T, pool *pgxpool.Pool, paymentID uint64) (paymententity.PaymentStatus, string) {
	t.Helper()

	var (
		status              paymententity.PaymentStatus
		providerReferenceID string
	)

	err := pool.QueryRow(
		context.Background(),
		`
			SELECT
				status,
				COALESCE(provider_reference_id, '')
			FROM payments
			WHERE id = $1
		`,
		paymentID,
	).Scan(
		&status,
		&providerReferenceID,
	)

	if err != nil {
		t.Fatalf("read payment state: %v", err)
	}
	return status, providerReferenceID
}
