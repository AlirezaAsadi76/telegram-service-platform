package postgrescheckouttesting

import (
	"context"
	"telegram-service-platform/entity/orderentity"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func getOrderState(t *testing.T, pool *pgxpool.Pool, orderID uint64,
) (orderentity.OrderStatus, *uint64, string) {
	t.Helper()

	var (
		status          orderentity.OrderStatus
		providerID      *uint64
		externalOrderID string
	)

	err := pool.QueryRow(
		context.Background(),
		`
		SELECT
			status,
			provider_id,
			external_order_id
		FROM orders
		WHERE id = $1
		`,
		orderID,
	).Scan(
		&status,
		&providerID,
		&externalOrderID,
	)

	if err != nil {
		t.Fatalf(
			"get order state: %v",
			err,
		)
	}

	return status, providerID, externalOrderID
}
