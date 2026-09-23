package postgresordertesting

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/repository/postgresorder"
)

func TestGetStaleProcessingWithoutEvidence_ReturnsOnlyEligibleOrders(
	t *testing.T,
) {
	pool := newTestPool(t)
	repo := postgresorder.NewWithExecutor(pool)

	staleAt := time.Now().Add(-1 * time.Hour)
	freshAt := time.Now()

	candidateID := createTestOrder(
		t,
		pool,
		orderentity.OrderStatusProcessing,
		ptr(""),
		staleAt,
	)

	externalIDOrderID := createTestOrder(
		t,
		pool,
		orderentity.OrderStatusProcessing,
		ptr("EXT-100"),
		staleAt,
	)

	freshOrderID := createTestOrder(
		t,
		pool,
		orderentity.OrderStatusProcessing,
		nil,
		freshAt,
	)

	attemptOrderID := createTestOrder(
		t,
		pool,
		orderentity.OrderStatusProcessing,
		nil,
		staleAt,
	)

	paidOrderID := createTestOrder(
		t,
		pool,
		orderentity.OrderStatusPaid,
		nil,
		staleAt,
	)

	createUnresolvedAttempt(
		t,
		pool,
		attemptOrderID,
	)

	t.Cleanup(func() {
		_, _ = pool.Exec(
			context.Background(),
			`DELETE FROM order_fulfillment_attempts WHERE order_id = ANY($1)`,
			[]uint64{attemptOrderID},
		)

		_, _ = pool.Exec(
			context.Background(),
			`DELETE FROM orders WHERE id = ANY($1)`,
			[]uint64{
				candidateID,
				externalIDOrderID,
				freshOrderID,
				attemptOrderID,
				paidOrderID,
			},
		)
	})

	orders, err := repo.GetStaleProcessingWithoutEvidence(
		context.Background(),
		30*time.Minute,
		50,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	foundCandidate := false

	for _, order := range orders {
		if order.ID == candidateID {
			foundCandidate = true
		}

		if order.ID == externalIDOrderID {
			t.Fatalf(
				"order with external ID %d must be excluded",
				externalIDOrderID,
			)
		}

		if order.ID == freshOrderID {
			t.Fatalf(
				"fresh PROCESSING order %d must be excluded",
				freshOrderID,
			)
		}

		if order.ID == attemptOrderID {
			t.Fatalf(
				"order with unresolved attempt %d must be excluded",
				attemptOrderID,
			)
		}

		if order.ID == paidOrderID {
			t.Fatalf(
				"PAID order %d must be excluded",
				paidOrderID,
			)
		}
	}

	if !foundCandidate {
		t.Fatalf(
			"expected candidate order ID %d in recovery result",
			candidateID,
		)
	}
}

func TestSaveExternalOrder_AcceptsInitiallyEmptyExternalOrderID(
	t *testing.T,
) {
	pool := newTestPool(t)
	repo := postgresorder.NewWithExecutor(pool)

	orderID := createTestOrder(
		t,
		pool,
		orderentity.OrderStatusProcessing,
		ptr(""),
		time.Now(),
	)

	t.Cleanup(func() {
		_, _ = pool.Exec(
			context.Background(),
			`DELETE FROM orders WHERE id = $1`,
			orderID,
		)
	})

	err := repo.SaveExternalOrder(
		context.Background(),
		orderID,
		7,
		"EXT-123",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var (
		providerID      sql.NullInt64
		externalOrderID string
	)

	err = pool.QueryRow(
		context.Background(),
		`
			SELECT provider_id, external_order_id
			FROM orders
			WHERE id = $1
		`,
		orderID,
	).Scan(
		&providerID,
		&externalOrderID,
	)
	if err != nil {
		t.Fatalf("read saved provider order: %v", err)
	}

	if !providerID.Valid {
		t.Fatal("expected provider ID to be persisted")
	}

	if providerID.Int64 != 7 {
		t.Fatalf(
			"expected provider ID 7, got %d",
			providerID.Int64,
		)
	}

	if externalOrderID != "EXT-123" {
		t.Fatalf(
			"expected external order ID EXT-123, got %s",
			externalOrderID,
		)
	}
}
