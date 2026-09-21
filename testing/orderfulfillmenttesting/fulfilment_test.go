package orderfulfillmenttesting

import (
	"context"
	"errors"
	"testing"
	"time"

	"telegram-service-platform/entity"
	"telegram-service-platform/entity/notificationentity"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/entity/productentity"
	"telegram-service-platform/entity/providerentity"
	"telegram-service-platform/params/smmparams"
	"telegram-service-platform/scheduler/jobs/orderfulfillerjob"
	"telegram-service-platform/scheduler/jobs/orderfulfillmentrecoveryjob"
	"telegram-service-platform/scheduler/jobs/statussyncJob"
	"telegram-service-platform/service/checkoutservice"
	"telegram-service-platform/service/notificationservice"
	"telegram-service-platform/service/orderfulfillmentservice"
	"telegram-service-platform/service/orderservice"
	"telegram-service-platform/service/smmproviderservice"
)

type fulfillmentE2ESystem struct {
	orderRepo *fakeOrderRepository
	queue     *fakeQueue

	providerRepo *fakeProviderRepository
	provider     *fakeSMMProvider

	notificationStore *fakeNotificationStore
	checkoutRepo      *fakeCheckoutTransactionRepository

	orderService        *orderservice.Service
	smmProviderService  *smmproviderservice.Service
	notificationService *notificationservice.Service
	checkoutService     *checkoutservice.Service
	fulfillmentEnqueuer *orderfulfillmentservice.Service

	orderFulfiller *orderfulfillerjob.Job
	statusSync     *statussyncjob.Job
	recovery       *orderfulfillmentrecoveryjob.Job
}

func newFulfillmentE2ESystem(
	order *orderentity.Order,
	provider *fakeSMMProvider,
	queueMessages [][]string,
	reconciliationAfter time.Duration,
) *fulfillmentE2ESystem {
	orderRepo := newFakeOrderRepository(order)

	queue := &fakeQueue{
		messages: queueMessages,
	}

	providerRepo := &fakeProviderRepository{
		provider: &providerentity.Provider{
			ID:       7,
			Name:     "provider-a",
			Type:     providerentity.ProviderTypeSMM,
			IsActive: true,
		},
	}

	notificationStore := &fakeNotificationStore{}

	checkoutRepo := &fakeCheckoutTransactionRepository{
		order: order,
	}

	orderService := orderservice.New(
		orderRepo,
	)

	smmProviderService := smmproviderservice.New(
		providerRepo,
		smmproviderservice.Config{
			FailureThreshold: 3,
			SuccessThreshold: 1,
		},
	)

	smmProviderService.RegisterProvider(
		"provider-a",
		provider,
	)

	notificationService := notificationservice.New(
		notificationStore,
		notificationStore,
		notificationservice.Config{},
	)

	checkoutService := checkoutservice.New(
		nil,
		nil,
		orderService,
		smmProviderService,
		notificationService,
		checkoutRepo,
		nil,
		nil,
		checkoutservice.Config{},
	)

	fulfillmentEnqueuer := orderfulfillmentservice.New(
		queue,
		orderfulfillmentservice.Config{
			QueueKey: "orders:paid",
		},
	)

	orderFulfiller := orderfulfillerjob.New(
		orderService,
		smmProviderService,
		notificationService,
		nil,
		checkoutService,
		queue,
		orderfulfillerjob.Config{
			QueueKey: "orders:paid",
			Timeout:  time.Second,
		},
	)

	statusSync := statussyncjob.New(
		orderService,
		smmProviderService,
		notificationService,
		checkoutService,
	)

	recovery := orderfulfillmentrecoveryjob.New(
		orderService,
		fulfillmentEnqueuer,
		orderfulfillmentrecoveryjob.Config{
			StaleAfter:          15 * time.Minute,
			ReconciliationAfter: reconciliationAfter,
			BatchSize:           50,
		},
	)

	return &fulfillmentE2ESystem{
		orderRepo: orderRepo,
		queue:     queue,

		providerRepo: providerRepo,
		provider:     provider,

		notificationStore: notificationStore,
		checkoutRepo:      checkoutRepo,

		orderService:        orderService,
		smmProviderService:  smmProviderService,
		notificationService: notificationService,
		checkoutService:     checkoutService,
		fulfillmentEnqueuer: fulfillmentEnqueuer,

		orderFulfiller: orderFulfiller,
		statusSync:     statusSync,
		recovery:       recovery,
	}
}

func newPaidSMMOrder() *orderentity.Order {
	return &orderentity.Order{
		ID:          42,
		UserID:      100,
		ProductType: productentity.ProductTypeSMM,
		ProductID:   55,
		Quantity:    1000,
		TargetLink:  "https://example.com/test",
		Status:      orderentity.OrderStatusPaid,
		Amount:      entity.Amount{},
		Currency:    entity.Currency(""),
		CreatedAt:   time.Now().Add(-20 * time.Minute),
	}
}

func TestFulfillmentE2E_CreatedThenCompleted(t *testing.T) {

	order := newPaidSMMOrder()

	provider := &fakeSMMProvider{
		createResponse: smmparams.CreateOrderAdapterResponse{
			Outcome:         smmparams.CreateOrderOutcomeCreated,
			ExternalOrderID: "EXT-100",
		},
		status: orderentity.OrderStatusProcessing,
	}

	system := newFulfillmentE2ESystem(
		order,
		provider,
		[][]string{
			{
				"orders:paid",
				"42",
			},
		},
		30*time.Minute,
	)

	// Act — fulfillment worker

	if err := system.orderFulfiller.Run(
		context.Background(),
	); err != nil {
		t.Fatalf(
			"unexpected fulfillment error: %v",
			err,
		)
	}

	// Assert — after fulfillment worker

	if order.Status != orderentity.OrderStatusProcessing {
		t.Fatalf(
			"expected PROCESSING after fulfillment, got %s",
			order.Status,
		)
	}

	if order.ProviderID == nil {
		t.Fatal("expected provider ID")
	}

	if *order.ProviderID != 7 {
		t.Fatalf(
			"expected provider ID 7, got %d",
			*order.ProviderID,
		)
	}

	if order.ExternalOrderID != "EXT-100" {
		t.Fatalf(
			"expected external order ID EXT-100, got %s",
			order.ExternalOrderID,
		)
	}

	if provider.createCalls != 1 {
		t.Fatalf(
			"expected one provider Create call, got %d",
			provider.createCalls,
		)
	}

	if len(system.orderRepo.attempts) != 1 {
		t.Fatalf(
			"expected one fulfillment attempt, got %d",
			len(system.orderRepo.attempts),
		)
	}

	if system.orderRepo.attempts[0].ResolvedAt == nil {
		t.Fatal("expected fulfillment attempt to be resolved")
	}

	// Change external provider state before StatusSync.

	provider.status = orderentity.OrderStatusCompleted

	// Act — status sync

	if err := system.statusSync.Run(
		context.Background(),
	); err != nil {
		t.Fatalf(
			"unexpected status sync error: %v",
			err,
		)
	}

	// Assert — final state

	if order.Status != orderentity.OrderStatusCompleted {
		t.Fatalf(
			"expected COMPLETED, got %s",
			order.Status,
		)
	}

	if provider.createCalls != 1 {
		t.Fatalf(
			"expected provider Create to remain 1, got %d",
			provider.createCalls,
		)
	}

	if provider.statusCalls != 1 {
		t.Fatalf(
			"expected one provider status call, got %d",
			provider.statusCalls,
		)
	}

	if len(system.notificationStore.created) != 2 {
		t.Fatalf(
			"expected processing and completed notifications, got %d",
			len(system.notificationStore.created),
		)
	}

	if system.notificationStore.created[0].Type !=
		notificationentity.NotificationTypeOrderProcessing {
		t.Fatalf(
			"expected first notification ORDER_PROCESSING, got %s",
			system.notificationStore.created[0].Type,
		)
	}

	if system.notificationStore.created[1].Type !=
		notificationentity.NotificationTypeOrderCompleted {
		t.Fatalf(
			"expected second notification ORDER_COMPLETED, got %s",
			system.notificationStore.created[1].Type,
		)
	}

	expectedOrderEvents := []string{
		"get_order",
		"claim",
		"save_external_order",
		"create_attempt",
		"resolve_attempt",
		"get_by_status",
		"complete_processing",
	}

	assertEventSequence(
		t,
		system.orderRepo.events,
		expectedOrderEvents,
	)
}

func TestFulfillmentE2E_RejectedThenRefunded(t *testing.T) {
	// Arrange

	order := newPaidSMMOrder()

	provider := &fakeSMMProvider{
		createResponse: smmparams.CreateOrderAdapterResponse{
			Outcome: smmparams.CreateOrderOutcomeRejected,
		},
	}

	system := newFulfillmentE2ESystem(
		order,
		provider,
		[][]string{
			{
				"orders:paid",
				"42",
			},
		},
		30*time.Minute,
	)

	// Act

	if err := system.orderFulfiller.Run(
		context.Background(),
	); err != nil {
		t.Fatalf(
			"unexpected fulfillment error: %v",
			err,
		)
	}

	// Assert

	if order.Status != orderentity.OrderStatusFailed {
		t.Fatalf(
			"expected FAILED after provider rejection refund, got %s",
			order.Status,
		)
	}

	if provider.createCalls != 1 {
		t.Fatalf(
			"expected one provider Create call, got %d",
			provider.createCalls,
		)
	}

	if len(system.checkoutRepo.calls) != 1 {
		t.Fatalf(
			"expected one refund call, got %d",
			len(system.checkoutRepo.calls),
		)
	}

	refundRequest := system.checkoutRepo.calls[0]

	if refundRequest.OrderID != order.ID {
		t.Fatalf(
			"expected refund order ID %d, got %d",
			order.ID,
			refundRequest.OrderID,
		)
	}

	if refundRequest.Reason != "smm_provider_rejected" {
		t.Fatalf(
			"expected reason smm_provider_rejected, got %s",
			refundRequest.Reason,
		)
	}

	if len(system.notificationStore.created) != 1 {
		t.Fatalf(
			"expected one failure notification, got %d",
			len(system.notificationStore.created),
		)
	}

	if system.notificationStore.created[0].Type !=
		notificationentity.NotificationTypeOrderFailed {
		t.Fatalf(
			"expected ORDER_FAILED notification, got %s",
			system.notificationStore.created[0].Type,
		)
	}

	// StatusSync must have nothing to do after refund.
	statusErr := system.statusSync.Run(
		context.Background(),
	)

	if statusErr != nil {
		t.Fatalf(
			"unexpected status sync error: %v",
			statusErr,
		)
	}

	if provider.statusCalls != 0 {
		t.Fatalf(
			"expected no provider status call after FAILED, got %d",
			provider.statusCalls,
		)
	}
}

func TestFulfillmentE2E_UnknownThenReconciliationRequired(t *testing.T) {
	// Arrange

	order := newPaidSMMOrder()

	provider := &fakeSMMProvider{
		createErr: errors.New("provider timeout"),
	}

	system := newFulfillmentE2ESystem(
		order,
		provider,
		[][]string{
			{
				"orders:paid",
				"42",
			},
		},
		30*time.Minute,
	)

	// Act — fulfillment worker

	if err := system.orderFulfiller.Run(
		context.Background(),
	); err != nil {
		t.Fatalf(
			"unexpected fulfillment error: %v",
			err,
		)
	}

	// Assert — immediately after UNKNOWN

	if order.Status != orderentity.OrderStatusProcessing {
		t.Fatalf(
			"expected PROCESSING, got %s",
			order.Status,
		)
	}

	if order.ProviderID == nil {
		t.Fatal("expected provider ID after UNKNOWN result")
	}

	if *order.ProviderID != 7 {
		t.Fatalf(
			"expected provider ID 7, got %d",
			*order.ProviderID,
		)
	}

	if order.ExternalOrderID != "" {
		t.Fatalf(
			"expected no external order ID for UNKNOWN result, got %s",
			order.ExternalOrderID,
		)
	}

	if len(system.orderRepo.attempts) != 1 {
		t.Fatalf(
			"expected one UNKNOWN fulfillment attempt, got %d",
			len(system.orderRepo.attempts),
		)
	}

	attempt := system.orderRepo.attempts[0]

	if attempt.Outcome !=
		orderentity.FulfillmentAttemptOutcomeUnknown {
		t.Fatalf(
			"expected UNKNOWN attempt, got %s",
			attempt.Outcome,
		)
	}

	if attempt.ResolvedAt != nil {
		t.Fatal(
			"expected UNKNOWN attempt to remain unresolved",
		)
	}

	// Make the unresolved attempt older than the reconciliation threshold.

	attempt.CreatedAt = time.Now().Add(-45 * time.Minute)

	// Act — recovery

	if err := system.recovery.Run(
		context.Background(),
	); err != nil {
		t.Fatalf(
			"unexpected recovery error: %v",
			err,
		)
	}

	// Assert — reconciliation state

	if order.Status != orderentity.OrderStatusProcessing {
		t.Fatalf(
			"expected PROCESSING after reconciliation signal, got %s",
			order.Status,
		)
	}

	if order.ExternalOrderID != "" {
		t.Fatalf(
			"expected external order ID to remain empty, got %s",
			order.ExternalOrderID,
		)
	}

	if system.orderRepo.attempts[0].ResolvedAt != nil {
		t.Fatal(
			"expected UNKNOWN attempt to remain unresolved",
		)
	}

	if provider.createCalls != 1 {
		t.Fatalf(
			"expected provider Create exactly once, got %d",
			provider.createCalls,
		)
	}

	if provider.statusCalls != 0 {
		t.Fatalf(
			"expected no provider status lookup without external ID, got %d",
			provider.statusCalls,
		)
	}

	if len(system.checkoutRepo.calls) != 0 {
		t.Fatalf(
			"expected no refund for UNKNOWN result, got %d",
			len(system.checkoutRepo.calls),
		)
	}
}

func TestFulfillmentE2E_DuplicateQueueMessageDoesNotDuplicateFulfillment(t *testing.T) {
	// Arrange

	order := newPaidSMMOrder()

	provider := &fakeSMMProvider{
		createResponse: smmparams.CreateOrderAdapterResponse{
			Outcome:         smmparams.CreateOrderOutcomeCreated,
			ExternalOrderID: "EXT-200",
		},
	}

	system := newFulfillmentE2ESystem(
		order,
		provider,
		[][]string{
			{
				"orders:paid",
				"42",
			},
			{
				"orders:paid",
				"42",
			},
		},
		30*time.Minute,
	)

	// Act — first delivery attempt

	if err := system.orderFulfiller.Run(
		context.Background(),
	); err != nil {
		t.Fatalf(
			"unexpected first fulfillment error: %v",
			err,
		)
	}

	// Act — duplicate queue message

	if err := system.orderFulfiller.Run(
		context.Background(),
	); err != nil {
		t.Fatalf(
			"unexpected duplicate fulfillment error: %v",
			err,
		)
	}

	// Assert

	if order.Status != orderentity.OrderStatusProcessing {
		t.Fatalf(
			"expected PROCESSING, got %s",
			order.Status,
		)
	}

	if provider.createCalls != 1 {
		t.Fatalf(
			"expected exactly one provider Create call, got %d",
			provider.createCalls,
		)
	}

	if len(system.orderRepo.attempts) != 1 {
		t.Fatalf(
			"expected exactly one fulfillment attempt, got %d",
			len(system.orderRepo.attempts),
		)
	}

	if system.orderRepo.claimCalls != 1 {
		t.Fatalf(
			"expected exactly one successful claim call, got %d",
			system.orderRepo.claimCalls,
		)
	}

	if len(system.notificationStore.created) != 1 {
		t.Fatalf(
			"expected exactly one processing notification, got %d",
			len(system.notificationStore.created),
		)
	}

	expectedEvents := []string{
		"get_order",
		"claim",
		"save_external_order",
		"create_attempt",
		"resolve_attempt",
		"get_order",
	}

	assertEventSequence(
		t,
		system.orderRepo.events,
		expectedEvents,
	)
}

func assertEventSequence(
	t *testing.T,
	actual []string,
	expected []string,
) {
	t.Helper()

	if len(actual) != len(expected) {
		t.Fatalf(
			"expected %d events, got %d: %#v",
			len(expected),
			len(actual),
			actual,
		)
	}

	for i, want := range expected {
		if actual[i] != want {
			t.Fatalf(
				"expected event %q at position %d, got %q",
				want,
				i,
				actual[i],
			)
		}
	}
}
