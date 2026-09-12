package postgrespaymenttesting

import (
	"context"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/pkg/richerror"
	"telegram-service-platform/repository/postgres"
	"telegram-service-platform/repository/postgrespayment"
	"testing"
)

func TestPaymentRepository_GetByExternalID_Success(t *testing.T) {
	pool := newTestPool(t)

	repo := postgrespayment.NewWithExecutor(
		pool,
		postgres.NewTransactionProvider(pool),
	)

	payment, _ := createTestCase(
		t,
		pool,
		paymententity.PaymentStatusPending,
		orderentity.OrderStatusPending,
	)

	_, err := pool.Exec(
		context.Background(),
		`
			UPDATE payments
			SET external_id = $1
			WHERE id = $2
		`,
		"AUTH-123",
		payment.ID,
	)

	if err != nil {
		t.Fatalf("set external ID: %v", err)
	}

	result, gErr := repo.GetByExternalID(
		context.Background(),
		"AUTH-123",
	)

	if gErr != nil {
		t.Fatalf("unexpected error: %v", gErr)
	}

	if result.ID != payment.ID {
		t.Fatalf("expected payment ID %d, got %d", payment.ID, result.ID)
	}

	if result.ExternalID != "AUTH-123" {
		t.Fatalf("expected external ID AUTH-123, got %s", result.ExternalID)
	}
}

func TestPaymentRepository_GetByExternalID_NotFound(t *testing.T) {
	pool := newTestPool(t)

	repo := postgrespayment.NewWithExecutor(
		pool,
		postgres.NewTransactionProvider(pool),
	)

	_, err := repo.GetByExternalID(
		context.Background(),
		"AUTH-NOT-FOUND",
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !richerror.IsCode(err, richerror.CodePaymentNotFound) {
		t.Fatalf("expected payment not found, got %v", err)
	}
}
