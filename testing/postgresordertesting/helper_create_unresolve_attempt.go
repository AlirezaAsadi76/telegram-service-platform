package postgresordertesting

import (
	"context"
	"telegram-service-platform/entity/orderentity"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func createUnresolvedAttempt(
	t *testing.T,
	pool *pgxpool.Pool,
	orderID uint64,
) {
	t.Helper()

	_, err := pool.Exec(
		context.Background(),
		`
			INSERT INTO order_fulfillment_attempts (
				order_id,
				provider_id,
				outcome,
				external_order_id,
				resolved_at,
				created_at,
				updated_at
			)
			VALUES ($1, $2, $3, $4, NULL, $5, $6)
		`,
		orderID,
		7,
		orderentity.FulfillmentAttemptOutcomeUnknown,
		nil,
		time.Now().Add(-45*time.Minute),
		time.Now().Add(-45*time.Minute),
	)
	if err != nil {
		t.Fatalf("create unresolved fulfillment attempt: %v", err)
	}
}

func ptr(value string) *string {
	return &value
}
