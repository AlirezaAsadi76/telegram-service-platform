package postgrespaymenttesting

import (
	"context"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/entity/paymententity"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

type testOrder struct {
	ID     uint64
	Status orderentity.OrderStatus
}

type testPayment struct {
	ID      uint64
	OrderID uint64
	Status  paymententity.PaymentStatus
}

func createTestOrder(t *testing.T, pool *pgxpool.Pool, status orderentity.OrderStatus) testOrder {
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
				amount,
				currency,
				status
			)
			VALUES ( 1, 'TEST', 1, 1, 100, 'IRR', $1)
			RETURNING id
		`,
		status,
	).Scan(&orderID)

	if err != nil {
		t.Fatalf("create test order: %v", err)
	}

	return testOrder{
		ID:     orderID,
		Status: status,
	}
}

func createTestPayment(t *testing.T, pool *pgxpool.Pool, orderID uint64, status paymententity.PaymentStatus) testPayment {
	t.Helper()

	var paymentID uint64

	err := pool.QueryRow(
		context.Background(),
		`
			INSERT INTO payments (
				order_id,
				user_id,
				method,
				amount,
				currency,
				status
			)
			VALUES ($1, 1, 'ZARINPAL', 100, 'IRR', $2)
			RETURNING id
		`,
		orderID, status).Scan(&paymentID)

	if err != nil {
		t.Fatalf("create test payment: %v", err)
	}

	return testPayment{
		ID:      paymentID,
		OrderID: orderID,
		Status:  status,
	}
}
