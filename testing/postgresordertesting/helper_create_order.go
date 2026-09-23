package postgresordertesting

import (
	"context"
	"telegram-service-platform/entity/orderentity"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func createTestOrder(
	t *testing.T,
	pool *pgxpool.Pool,
	status orderentity.OrderStatus,
	externalOrderID *string,
	updatedAt time.Time,
) uint64 {
	t.Helper()

	var orderID uint64

	err := pool.QueryRow(
		context.Background(),
		`
			INSERT INTO orders (
				user_id,
				product_type,
				product_id,
				quantity,
				target_link,
				amount,
				currency,
				status,
				external_order_id,
				provider_id,
				metadata,
				created_at,
				updated_at
			)
			VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
			)
			RETURNING id
		`,
		100,
		"SMM",
		1,
		1,
		"https://example.com",
		10,
		"TOMAN",
		status,
		externalOrderID,
		nil,
		nil,
		updatedAt,
		updatedAt,
	).Scan(&orderID)
	if err != nil {
		t.Fatalf("create test order: %v", err)
	}

	return orderID
}
