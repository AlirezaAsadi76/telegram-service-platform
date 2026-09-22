package orderfulfillmentrecoverytesting

import (
	"context"
	"errors"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/scheduler/jobs/orderfulfillmentrecoveryjob"
	"telegram-service-platform/service/orderfulfillmentservice"
	"telegram-service-platform/service/orderservice"
	"testing"
	"time"
)

func newRecoveryJob(repo *fakeOrderRepository, queue *fakeQueue, config orderfulfillmentrecoveryjob.Config) *orderfulfillmentrecoveryjob.Job {
	orderService := orderservice.New(repo)

	fulfillmentService := orderfulfillmentservice.New(
		queue,
		orderfulfillmentservice.Config{
			QueueKey: "orders:paid",
		},
	)

	return orderfulfillmentrecoveryjob.New(
		orderService,
		fulfillmentService,
		config,
	)
}

func TestJob_Run_RecoverCreatedAttempt(t *testing.T) {

	order := &orderentity.Order{
		ID:     42,
		UserID: 100,
		Status: orderentity.OrderStatusProcessing,
	}

	attempt := orderentity.FulfillmentAttempt{
		ID:              10,
		OrderID:         42,
		ProviderID:      7,
		Outcome:         orderentity.FulfillmentAttemptOutcomeCreated,
		ExternalOrderID: "EXT-123",
		CreatedAt:       time.Now().Add(-5 * time.Minute),
	}

	repo := &fakeOrderRepository{
		order:              order,
		unresolvedAttempts: []orderentity.FulfillmentAttempt{attempt},
	}

	queue := &fakeQueue{}

	job := newRecoveryJob(
		repo,
		queue,
		orderfulfillmentrecoveryjob.Config{
			ReconciliationAfter: 30 * time.Minute,
			BatchSize:           50,
		},
	)

	err := job.Run(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if order.ProviderID == nil {
		t.Fatal("expected provider ID to be recovered")
	}

	if *order.ProviderID != 7 {
		t.Fatalf(
			"expected provider ID 7, got %d",
			*order.ProviderID,
		)
	}

	if order.ExternalOrderID != "EXT-123" {
		t.Fatalf(
			"expected external order ID EXT-123, got %s",
			order.ExternalOrderID,
		)
	}

	if repo.saveExternalCalls != 1 {
		t.Fatalf(
			"expected SaveExternalOrder once, got %d",
			repo.saveExternalCalls,
		)
	}

	if repo.resolveCalls != 1 {
		t.Fatalf(
			"expected Resolve once, got %d",
			repo.resolveCalls,
		)
	}

	if repo.unresolvedAttempts[0].ResolvedAt == nil {
		t.Fatal("expected attempt to be resolved")
	}

	expectedEvents := []string{
		"get_unresolved_attempts",
		"get_order",
		"save_external_order",
		"resolve_attempt",
		"get_stale_paid",
	}

	for i, expected := range expectedEvents {
		if i >= len(repo.events) {
			t.Fatalf(
				"missing event %q",
				expected,
			)
		}

		if repo.events[i] != expected {
			t.Fatalf(
				"expected event %q at position %d, got %q",
				expected,
				i,
				repo.events[i],
			)
		}
	}
}

func TestJob_Run_ResolvesAlreadyPersistedCreatedAttempt(t *testing.T) {

	providerID := uint64(7)

	order := &orderentity.Order{
		ID:              42,
		Status:          orderentity.OrderStatusProcessing,
		ProviderID:      &providerID,
		ExternalOrderID: "EXT-123",
	}

	attempt := orderentity.FulfillmentAttempt{
		ID:              10,
		OrderID:         42,
		ProviderID:      7,
		Outcome:         orderentity.FulfillmentAttemptOutcomeCreated,
		ExternalOrderID: "EXT-123",
		CreatedAt:       time.Now().Add(-5 * time.Minute),
	}

	repo := &fakeOrderRepository{
		order:              order,
		unresolvedAttempts: []orderentity.FulfillmentAttempt{attempt},
	}

	queue := &fakeQueue{}

	job := newRecoveryJob(
		repo,
		queue,
		orderfulfillmentrecoveryjob.Config{
			BatchSize: 50,
		},
	)

	err := job.Run(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.saveExternalCalls != 0 {
		t.Fatalf(
			"expected SaveExternalOrder not to be called, got %d",
			repo.saveExternalCalls,
		)
	}

	if repo.resolveCalls != 1 {
		t.Fatalf(
			"expected Resolve once, got %d",
			repo.resolveCalls,
		)
	}

	if repo.unresolvedAttempts[0].ResolvedAt == nil {
		t.Fatal("expected attempt to be resolved")
	}

	expectedEvents := []string{
		"get_unresolved_attempts",
		"get_order",
		"resolve_attempt",
		"get_stale_paid",
	}

	for i, expected := range expectedEvents {
		if i >= len(repo.events) {
			t.Fatalf("missing event %q", expected)
		}

		if repo.events[i] != expected {
			t.Fatalf(
				"expected event %q at position %d, got %q",
				expected,
				i,
				repo.events[i],
			)
		}
	}
}

func TestJob_Run_UnknownAttemptRemainsPendingBeforeThreshold(t *testing.T) {

	providerID := uint64(7)

	order := &orderentity.Order{
		ID:         42,
		Status:     orderentity.OrderStatusProcessing,
		ProviderID: &providerID,
	}

	attempt := orderentity.FulfillmentAttempt{
		ID:         10,
		OrderID:    42,
		ProviderID: 7,
		Outcome:    orderentity.FulfillmentAttemptOutcomeUnknown,
		CreatedAt:  time.Now().Add(-5 * time.Minute),
	}

	repo := &fakeOrderRepository{
		order:              order,
		unresolvedAttempts: []orderentity.FulfillmentAttempt{attempt},
	}

	queue := &fakeQueue{}

	job := newRecoveryJob(
		repo,
		queue,
		orderfulfillmentrecoveryjob.Config{
			ReconciliationAfter: 30 * time.Minute,
			BatchSize:           50,
		},
	)

	err := job.Run(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.assignProviderCalls != 0 {
		t.Fatalf(
			"expected AssignProvider not to be called, got %d",
			repo.assignProviderCalls,
		)
	}

	if repo.saveExternalCalls != 0 {
		t.Fatalf(
			"expected SaveExternalOrder not to be called, got %d",
			repo.saveExternalCalls,
		)
	}

	if repo.resolveCalls != 0 {
		t.Fatalf(
			"expected Resolve not to be called, got %d",
			repo.resolveCalls,
		)
	}

	if order.Status != orderentity.OrderStatusProcessing {
		t.Fatalf(
			"expected PROCESSING, got %s",
			order.Status,
		)
	}

	if repo.unresolvedAttempts[0].ResolvedAt != nil {
		t.Fatal("expected attempt to remain unresolved")
	}
}

func TestJob_Run_UnknownAttemptRequiresReconciliationAfterThreshold(t *testing.T) {

	providerID := uint64(7)

	order := &orderentity.Order{
		ID:         42,
		Status:     orderentity.OrderStatusProcessing,
		ProviderID: &providerID,
	}

	attempt := orderentity.FulfillmentAttempt{
		ID:         10,
		OrderID:    42,
		ProviderID: 7,
		Outcome:    orderentity.FulfillmentAttemptOutcomeUnknown,
		CreatedAt:  time.Now().Add(-45 * time.Minute),
	}

	repo := &fakeOrderRepository{
		order:              order,
		unresolvedAttempts: []orderentity.FulfillmentAttempt{attempt},
	}

	queue := &fakeQueue{}

	job := newRecoveryJob(
		repo,
		queue,
		orderfulfillmentrecoveryjob.Config{
			ReconciliationAfter: 30 * time.Minute,
			BatchSize:           50,
		},
	)
	err := job.Run(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.assignProviderCalls != 0 {
		t.Fatalf(
			"expected AssignProvider not to be called, got %d",
			repo.assignProviderCalls,
		)
	}

	if repo.resolveCalls != 0 {
		t.Fatalf(
			"expected Resolve not to be called, got %d",
			repo.resolveCalls,
		)
	}

	if order.Status != orderentity.OrderStatusProcessing {
		t.Fatalf(
			"expected PROCESSING, got %s",
			order.Status,
		)
	}

	if repo.unresolvedAttempts[0].ResolvedAt != nil {
		t.Fatal("expected attempt to remain unresolved")
	}
}

func TestJob_Run_UnknownAttemptRecoversProviderAssignment(t *testing.T) {

	order := &orderentity.Order{
		ID:     42,
		Status: orderentity.OrderStatusProcessing,
	}

	attempt := orderentity.FulfillmentAttempt{
		ID:         10,
		OrderID:    42,
		ProviderID: 7,
		Outcome:    orderentity.FulfillmentAttemptOutcomeUnknown,
		CreatedAt:  time.Now().Add(-45 * time.Minute),
	}

	repo := &fakeOrderRepository{
		order:              order,
		unresolvedAttempts: []orderentity.FulfillmentAttempt{attempt},
	}

	queue := &fakeQueue{}

	job := newRecoveryJob(
		repo,
		queue,
		orderfulfillmentrecoveryjob.Config{
			ReconciliationAfter: 30 * time.Minute,
			BatchSize:           50,
		},
	)

	err := job.Run(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if order.ProviderID == nil {
		t.Fatal("expected provider ID to be assigned")
	}

	if *order.ProviderID != 7 {
		t.Fatalf(
			"expected provider ID 7, got %d",
			*order.ProviderID,
		)
	}

	if repo.assignProviderCalls != 1 {
		t.Fatalf(
			"expected AssignProvider once, got %d",
			repo.assignProviderCalls,
		)
	}

	if repo.resolveCalls != 0 {
		t.Fatalf(
			"expected attempt not to be resolved, got %d",
			repo.resolveCalls,
		)
	}

	if repo.saveExternalCalls != 0 {
		t.Fatalf(
			"expected SaveExternalOrder not to be called, got %d",
			repo.saveExternalCalls,
		)
	}
}

func TestJob_Run_ReEnqueuesStalePaidOrders(t *testing.T) {

	orders := []*orderentity.Order{
		{
			ID:     42,
			Status: orderentity.OrderStatusPaid,
		},
		{
			ID:     43,
			Status: orderentity.OrderStatusPaid,
		},
	}

	repo := &fakeOrderRepository{
		staleOrders: orders,
	}

	queue := &fakeQueue{}

	job := newRecoveryJob(
		repo,
		queue,
		orderfulfillmentrecoveryjob.Config{
			StaleAfter: 15 * time.Minute,
			BatchSize:  50,
		},
	)

	err := job.Run(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(queue.values) != 2 {
		t.Fatalf(
			"expected 2 enqueue operations, got %d",
			len(queue.values),
		)
	}

	if queue.queueKey != "orders:paid" {
		t.Fatalf(
			"expected queue key orders:paid, got %s",
			queue.queueKey,
		)
	}

	if queue.values[0] != uint64(42) {
		t.Fatalf(
			"expected first order ID 42, got %v",
			queue.values[0],
		)
	}

	if queue.values[1] != uint64(43) {
		t.Fatalf(
			"expected second order ID 43, got %v",
			queue.values[1],
		)
	}
}

func TestJob_Run_NoStalePaidOrdersDoesNotEnqueue(t *testing.T) {

	repo := &fakeOrderRepository{
		staleOrders: nil,
	}

	queue := &fakeQueue{}

	job := newRecoveryJob(
		repo,
		queue,
		orderfulfillmentrecoveryjob.Config{
			StaleAfter: 15 * time.Minute,
			BatchSize:  50,
		},
	)

	err := job.Run(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(queue.values) != 0 {
		t.Fatalf(
			"expected no enqueue operation, got %d",
			len(queue.values),
		)
	}

	expectedEvents := []string{
		"get_unresolved_attempts",
		"get_stale_paid",
	}

	for i, expected := range expectedEvents {
		if i >= len(repo.events) {
			t.Fatalf("missing event %q", expected)
		}

		if repo.events[i] != expected {
			t.Fatalf(
				"expected event %q at position %d, got %q",
				expected,
				i,
				repo.events[i],
			)
		}
	}
}

func TestJob_Run_PartialReenqueueFailureContinues(t *testing.T) {

	orders := []*orderentity.Order{
		{
			ID:     42,
			Status: orderentity.OrderStatusPaid,
		},
		{
			ID:     43,
			Status: orderentity.OrderStatusPaid,
		},
		{
			ID:     44,
			Status: orderentity.OrderStatusPaid,
		},
	}

	repo := &fakeOrderRepository{
		staleOrders: orders,
	}

	queue := &fakeQueue{
		failOrderIDs: map[uint64]error{
			43: errors.New("redis unavailable"),
		},
	}

	job := newRecoveryJob(
		repo,
		queue,
		orderfulfillmentrecoveryjob.Config{
			StaleAfter: 15 * time.Minute,
			BatchSize:  50,
		},
	)

	err := job.Run(context.Background())

	if err != nil {
		t.Fatalf(
			"expected partial enqueue failure to be handled, got %v",
			err,
		)
	}

	if len(queue.values) != 3 {
		t.Fatalf(
			"expected 3 enqueue attempts, got %d",
			len(queue.values),
		)
	}

	expected := []uint64{
		42,
		43,
		44,
	}

	for i, value := range queue.values {
		orderID, ok := value.(uint64)
		if !ok {
			t.Fatalf(
				"expected uint64 queue value, got %T",
				value,
			)
		}

		if orderID != expected[i] {
			t.Fatalf(
				"expected order ID %d at position %d, got %d",
				expected[i],
				i,
				orderID,
			)
		}
	}
}

func TestJob_Run_GetStalePaidFailure(t *testing.T) {

	repo := &fakeOrderRepository{
		stalePaidErr: errors.New("database unavailable"),
	}

	queue := &fakeQueue{}

	job := newRecoveryJob(
		repo,
		queue,
		orderfulfillmentrecoveryjob.Config{
			StaleAfter: 15 * time.Minute,
			BatchSize:  50,
		},
	)

	err := job.Run(context.Background())

	if err == nil {
		t.Fatal("expected error")
	}

	if len(queue.values) != 0 {
		t.Fatalf(
			"expected no enqueue after GetStalePaid failure, got %d",
			len(queue.values),
		)
	}
}

func TestJob_Run_StaleProcessingWithoutEvidenceRequiresReconciliation(
	t *testing.T,
) {
	order := &orderentity.Order{
		ID:              42,
		UserID:          100,
		Status:          orderentity.OrderStatusProcessing,
		ExternalOrderID: "",
		UpdatedAt:       time.Now().Add(-45 * time.Minute),
	}

	repo := &fakeOrderRepository{
		staleProcessingOrders: []*orderentity.Order{
			order,
		},
	}

	queue := &fakeQueue{}

	job := newRecoveryJob(
		repo,
		queue,
		orderfulfillmentrecoveryjob.Config{
			StaleAfter:          15 * time.Minute,
			ReconciliationAfter: 30 * time.Minute,
			BatchSize:           50,
		},
	)

	err := job.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.staleProcessingCalls != 1 {
		t.Fatalf(
			"expected GetStaleProcessingWithoutEvidence once, got %d",
			repo.staleProcessingCalls,
		)
	}

	if order.Status != orderentity.OrderStatusProcessing {
		t.Fatalf(
			"expected order to remain PROCESSING, got %s",
			order.Status,
		)
	}

	if order.ExternalOrderID != "" {
		t.Fatalf(
			"expected external order ID to remain empty, got %s",
			order.ExternalOrderID,
		)
	}

	if len(queue.values) != 0 {
		t.Fatalf(
			"expected no queue operation for stale PROCESSING order, got %d",
			len(queue.values),
		)
	}
}
