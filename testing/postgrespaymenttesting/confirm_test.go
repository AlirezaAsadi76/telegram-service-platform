package postgrespaymenttesting

import (
	"context"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/pkg/richerror"
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

func TestPaymentConfirmationRepository_Confirm_PaymentNotFound(t *testing.T) {
	pool := newTestPool(t)

	repo := postgrespayment.NewWithExecutor(
		pool,
		postgres.NewTransactionProvider(pool),
	)

	err := repo.Confirm(
		context.Background(),
		999999999,
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !richerror.IsCode(
		err,
		richerror.CodePaymentNotFound,
	) {
		t.Fatalf(
			"expected code %s, got %v",
			richerror.CodePaymentNotFound,
			err,
		)
	}
}

func TestPaymentConfirmationRepository_Confirm_AlreadyConfirmed(t *testing.T) {
	pool := newTestPool(t)

	repo := postgrespayment.NewWithExecutor(
		pool,
		postgres.NewTransactionProvider(pool),
	)

	payment, _ := createTestCase(
		t,
		pool,
		paymententity.PaymentStatusSuccess,
		orderentity.OrderStatusPending,
	)

	err := repo.Confirm(
		context.Background(),
		payment.ID,
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !richerror.IsCode(
		err,
		richerror.CodePaymentAlreadyConfirmed,
	) {
		t.Fatalf(
			"expected code %s, got %v",
			richerror.CodePaymentAlreadyConfirmed,
			err,
		)
	}
}

func TestPaymentConfirmationRepository_Confirm_Concurrent(t *testing.T) {
	pool := newTestPool(t)

	repo := postgrespayment.NewWithExecutor(
		pool,
		postgres.NewTransactionProvider(pool),
	)

	payment, order := createTestCase(t, pool,
		paymententity.PaymentStatusPending,
		orderentity.OrderStatusPending,
	)

	const countCallback = 5

	type result struct {
		err error
	}

	results := make(chan result, countCallback)

	ctx := context.Background()

	for i := 0; i < countCallback; i++ {
		go func() {
			results <- result{
				err: repo.Confirm(ctx, payment.ID),
			}
		}()

	}

	resResults := make([]result, 0, countCallback)

	for i := 0; i < countCallback; i++ {
		resResults = append(resResults, <-results)
	}

	successCount := 0
	alreadyConfirmedCount := 0

	for _, val := range resResults {
		if val.err == nil {
			successCount++
			continue
		}

		if richerror.IsCode(
			val.err,
			richerror.CodePaymentAlreadyConfirmed,
		) {
			alreadyConfirmedCount++
			continue
		}

		t.Fatalf(
			"unexpected error: %v",
			val.err,
		)
	}

	if successCount != 1 {
		t.Fatalf(
			"expected 1 successful confirmation, got %d",
			successCount,
		)
	}

	if alreadyConfirmedCount != countCallback-1 {
		t.Fatalf(
			"expected %d already-confirmed result, got %d",
			countCallback-1,
			alreadyConfirmedCount,
		)
	}

	var paymentStatus paymententity.PaymentStatus

	err := pool.QueryRow(
		context.Background(),
		`SELECT status FROM payments WHERE id = $1`,
		payment.ID,
	).Scan(&paymentStatus)

	if err != nil {
		t.Fatalf("read payment: %v", err)
	}

	if paymentStatus != paymententity.PaymentStatusSuccess {
		t.Fatalf(
			"expected payment SUCCESS, got %s",
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
		t.Fatalf("read order: %v", err)
	}

	if orderStatus != orderentity.OrderStatusPaid {
		t.Fatalf(
			"expected order PAID, got %s",
			orderStatus,
		)
	}
}
