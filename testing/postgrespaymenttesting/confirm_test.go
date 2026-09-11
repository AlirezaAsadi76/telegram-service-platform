package postgrespaymenttesting

import (
	"context"
	"telegram-service-platform/entity/orderentity"
	"testing"

	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/repository/postgres"
	"telegram-service-platform/repository/postgrespayment"
)

func TestPaymentConfirmationRepository_Confirm_Success(t *testing.T) {
	pool := newTestPool(t)

	repo := postgrespayment.NewWithExecutor(
		pool,
		postgres.NewTransactionProvider(pool),
	)

	payment, order := createTestCase(t, pool,
		paymententity.PaymentStatusPending,
		orderentity.OrderStatusPending,
	)

	err := repo.Confirm(
		context.Background(),
		payment.ID,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var paymentStatus paymententity.PaymentStatus

	err = pool.QueryRow(
		context.Background(),
		`SELECT status FROM payments WHERE id = $1`,
		payment.ID,
	).Scan(&paymentStatus)

	if err != nil {
		t.Fatalf("read payment status: %v", err)
	}

	if paymentStatus != paymententity.PaymentStatusSuccess {
		t.Fatalf(
			"expected payment status %s, got %s",
			paymententity.PaymentStatusSuccess,
			paymentStatus,
		)
	}

	var orderStatus orderentity.OrderStatus

	err = pool.QueryRow(
		context.Background(),
		`SELECT status FROM orders WHERE id = $1`,
		order.ID,
	).Scan(&orderStatus)

	if err != nil {
		t.Fatalf("read order status: %v", err)
	}

	if orderStatus != orderentity.OrderStatusPaid {
		t.Fatalf(
			"expected order status PAID, got %s",
			orderStatus,
		)
	}
}

func TestPaymentConfirmationRepository_Confirm_RollbackOnOrderFailure(t *testing.T) {
	pool := newTestPool(t)

	repo := postgrespayment.NewWithExecutor(
		pool,
		postgres.NewTransactionProvider(pool),
	)

	payment, order := createTestCase(
		t,
		pool,
		paymententity.PaymentStatusPending,
		orderentity.OrderStatusCanceled,
	)

	err := repo.Confirm(context.Background(), payment.ID)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var paymentStatus paymententity.PaymentStatus

	err = pool.QueryRow(
		context.Background(),
		`SELECT status FROM payments WHERE id = $1`,
		payment.ID,
	).Scan(&paymentStatus)

	if err != nil {
		t.Fatalf("read payment status: %v", err)
	}

	if paymentStatus != paymententity.PaymentStatusPending {
		t.Fatalf(
			"expected payment status PENDING after rollback, got %s",
			paymentStatus,
		)
	}

	var orderStatus orderentity.OrderStatus

	err = pool.QueryRow(
		context.Background(),
		`SELECT status FROM orders WHERE id = $1`,
		order.ID,
	).Scan(&orderStatus)

	if err != nil {
		t.Fatalf("read order status: %v", err)
	}

	if orderStatus != orderentity.OrderStatusCanceled {
		t.Fatalf(
			"expected order status CANCELED, got %s",
			orderStatus,
		)
	}
}
