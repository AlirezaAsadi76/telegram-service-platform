package paymentflowtesting

import (
	"context"
	"telegram-service-platform/entity/orderentity"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func readOrderStatus(t *testing.T, pool *pgxpool.Pool, orderID uint64) orderentity.OrderStatus {
	t.Helper()

	var status orderentity.OrderStatus

	err := pool.QueryRow(
		context.Background(),
		`
			SELECT status
			FROM orders
			WHERE id = $1
		`,
		orderID,
	).Scan(&status)

	if err != nil {
		t.Fatalf("read order status: %v", err)
	}

	return status
}
