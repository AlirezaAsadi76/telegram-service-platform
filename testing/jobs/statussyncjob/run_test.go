package statussynctesting

import (
	"context"
	"errors"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/scheduler/jobs/statussyncJob"
	"telegram-service-platform/service/checkoutservice"
	"telegram-service-platform/service/notificationservice"
	"telegram-service-platform/service/orderservice"
	"telegram-service-platform/service/smmproviderservice"
	"testing"
)

func newStatusSyncJob(
	orderRepo *fakeOrderRepository,
	providerRepo *fakeProviderRepository,
	provider *fakeSMMProvider,
	checkoutRepo *fakeCheckoutTransactionRepository,
) (
	*statussyncjob.Job,
	*fakeNotificationRepository,
	*fakeNotificationRedis,
) {
	orderService := orderservice.New(orderRepo)

	smmService := smmproviderservice.New(
		providerRepo,
		smmproviderservice.Config{
			FailureThreshold: 3,
			SuccessThreshold: 1,
		},
	)

	smmService.RegisterProvider(
		providerRepo.provider.Name,
		provider,
	)

	notificationRepo := &fakeNotificationRepository{}
	notificationRedis := &fakeNotificationRedis{}

	notificationService := notificationservice.New(
		notificationRepo,
		notificationRedis,
		notificationservice.Config{},
	)

	checkoutService := checkoutservice.New(
		nil,
		nil,
		orderService,
		smmService,
		notificationService,
		checkoutRepo,
		nil,
		nil,
		checkoutservice.Config{},
	)

	job := statussyncjob.New(
		orderService,
		smmService,
		notificationService,
		checkoutService,
	)

	return job, notificationRepo, notificationRedis
}

func TestJob_Run_CompletedOrder(t *testing.T) {

	providerID := uint64(7)

	order := &orderentity.Order{
		ID:              42,
		UserID:          100,
		Status:          orderentity.OrderStatusProcessing,
		ProviderID:      &providerID,
		ExternalOrderID: "EXT-100",
	}

	orderRepo := &fakeOrderRepository{
		orders:         []*orderentity.Order{order},
		completeResult: true,
	}

	providerRepo := &fakeProviderRepository{
		provider: newProvider(7, "provider-a"),
	}

	provider := &fakeSMMProvider{
		status: orderentity.OrderStatusCompleted,
	}

	checkoutRepo := &fakeCheckoutTransactionRepository{}

	job, notificationRepo, notificationRedis := newStatusSyncJob(
		orderRepo,
		providerRepo,
		provider,
		checkoutRepo,
	)

	err := job.Run(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if orderRepo.requestedStatus != orderentity.OrderStatusProcessing {
		t.Fatalf(
			"expected PROCESSING query, got %s",
			orderRepo.requestedStatus,
		)
	}

	if provider.statusCalls != 1 {
		t.Fatalf(
			"expected one status call, got %d",
			provider.statusCalls,
		)
	}

	if order.Status != orderentity.OrderStatusCompleted {
		t.Fatalf(
			"expected COMPLETED, got %s",
			order.Status,
		)
	}

	if len(notificationRepo.created) != 1 {
		t.Fatalf(
			"expected one notification, got %d",
			len(notificationRepo.created),
		)
	}

	if notificationRepo.created[0].UserID != 100 {
		t.Fatalf(
			"expected notification user ID 100, got %d",
			notificationRepo.created[0].UserID,
		)
	}

	if len(notificationRedis.enqueued) != 1 {
		t.Fatalf(
			"expected one notification enqueue, got %d",
			len(notificationRedis.enqueued),
		)
	}
}

func TestJob_Run_ProcessingOrderDoesNothing(t *testing.T) {

	providerID := uint64(7)

	order := &orderentity.Order{
		ID:              42,
		UserID:          100,
		Status:          orderentity.OrderStatusProcessing,
		ProviderID:      &providerID,
		ExternalOrderID: "EXT-100",
	}

	orderRepo := &fakeOrderRepository{
		orders: []*orderentity.Order{order},
	}

	providerRepo := &fakeProviderRepository{
		provider: newProvider(7, "provider-a"),
	}

	provider := &fakeSMMProvider{
		status: orderentity.OrderStatusProcessing,
	}

	checkoutRepo := &fakeCheckoutTransactionRepository{}

	job, notificationRepo, _ := newStatusSyncJob(
		orderRepo,
		providerRepo,
		provider,
		checkoutRepo,
	)

	err := job.Run(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if provider.statusCalls != 1 {
		t.Fatalf(
			"expected one status call, got %d",
			provider.statusCalls,
		)
	}

	if order.Status != orderentity.OrderStatusProcessing {
		t.Fatalf(
			"expected PROCESSING, got %s",
			order.Status,
		)
	}

	if orderRepo.events[len(orderRepo.events)-1] == "complete_processing" {
		t.Fatal("CompleteProcessing must not be called")
	}

	if len(notificationRepo.created) != 0 {
		t.Fatalf(
			"expected no notification, got %d",
			len(notificationRepo.created),
		)
	}
}

func TestJob_Run_FailedOrderRefundsThroughCheckout(t *testing.T) {

	providerID := uint64(7)

	order := &orderentity.Order{
		ID:              42,
		UserID:          100,
		Status:          orderentity.OrderStatusProcessing,
		ProviderID:      &providerID,
		ExternalOrderID: "EXT-100",
	}

	orderRepo := &fakeOrderRepository{
		orders: []*orderentity.Order{order},
	}

	providerRepo := &fakeProviderRepository{
		provider: newProvider(7, "provider-a"),
	}

	provider := &fakeSMMProvider{
		status: orderentity.OrderStatusFailed,
	}

	checkoutRepo := &fakeCheckoutTransactionRepository{}

	job, notificationRepo, _ := newStatusSyncJob(
		orderRepo,
		providerRepo,
		provider,
		checkoutRepo,
	)

	err := job.Run(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if provider.statusCalls != 1 {
		t.Fatalf(
			"expected one status call, got %d",
			provider.statusCalls,
		)
	}

	if len(checkoutRepo.calls) != 1 {
		t.Fatalf(
			"expected one refund call, got %d",
			len(checkoutRepo.calls),
		)
	}

	refundReq := checkoutRepo.calls[0]

	if refundReq.OrderID != 42 {
		t.Fatalf(
			"expected refund order ID 42, got %d",
			refundReq.OrderID,
		)
	}

	if refundReq.Reason != "provider_failed" {
		t.Fatalf(
			"expected reason provider_failed, got %s",
			refundReq.Reason,
		)
	}

	if len(notificationRepo.created) != 1 {
		t.Fatalf(
			"expected one failure notification, got %d",
			len(notificationRepo.created),
		)
	}
}

func TestJob_Run_SkipsOrderWithoutProviderData(t *testing.T) {

	order := &orderentity.Order{
		ID:         42,
		UserID:     100,
		Status:     orderentity.OrderStatusProcessing,
		ProviderID: nil,
	}

	orderRepo := &fakeOrderRepository{
		orders: []*orderentity.Order{order},
	}

	providerRepo := &fakeProviderRepository{
		provider: newProvider(7, "provider-a"),
	}

	provider := &fakeSMMProvider{
		status: orderentity.OrderStatusCompleted,
	}

	checkoutRepo := &fakeCheckoutTransactionRepository{}

	job, notificationRepo, _ := newStatusSyncJob(
		orderRepo,
		providerRepo,
		provider,
		checkoutRepo,
	)

	err := job.Run(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if provider.statusCalls != 0 {
		t.Fatalf(
			"expected no provider status call, got %d",
			provider.statusCalls,
		)
	}

	if order.Status != orderentity.OrderStatusProcessing {
		t.Fatalf(
			"expected PROCESSING, got %s",
			order.Status,
		)
	}

	if len(notificationRepo.created) != 0 {
		t.Fatalf(
			"expected no notification, got %d",
			len(notificationRepo.created),
		)
	}

	if len(checkoutRepo.calls) != 0 {
		t.Fatalf(
			"expected no refund, got %d",
			len(checkoutRepo.calls),
		)
	}
}

func TestJob_Run_ProviderStatusErrorDoesNotChangeOrder(t *testing.T) {

	providerID := uint64(7)

	order := &orderentity.Order{
		ID:              42,
		UserID:          100,
		Status:          orderentity.OrderStatusProcessing,
		ProviderID:      &providerID,
		ExternalOrderID: "EXT-100",
	}

	orderRepo := &fakeOrderRepository{
		orders: []*orderentity.Order{order},
	}

	providerRepo := &fakeProviderRepository{
		provider: newProvider(7, "provider-a"),
	}

	provider := &fakeSMMProvider{
		err: errors.New("provider timeout"),
	}

	checkoutRepo := &fakeCheckoutTransactionRepository{}

	job, notificationRepo, _ := newStatusSyncJob(
		orderRepo,
		providerRepo,
		provider,
		checkoutRepo,
	)

	err := job.Run(context.Background())

	if err != nil {
		t.Fatalf("expected worker to continue, got %v", err)
	}

	if provider.statusCalls != 1 {
		t.Fatalf(
			"expected one provider call, got %d",
			provider.statusCalls,
		)
	}

	if order.Status != orderentity.OrderStatusProcessing {
		t.Fatalf(
			"expected PROCESSING, got %s",
			order.Status,
		)
	}

	if len(checkoutRepo.calls) != 0 {
		t.Fatalf(
			"expected no refund, got %d",
			len(checkoutRepo.calls),
		)
	}

	if len(notificationRepo.created) != 0 {
		t.Fatalf(
			"expected no notification, got %d",
			len(notificationRepo.created),
		)
	}
}

func TestJob_Run_CompletionRaceDoesNotNotify(t *testing.T) {

	providerID := uint64(7)

	order := &orderentity.Order{
		ID:              42,
		UserID:          100,
		Status:          orderentity.OrderStatusProcessing,
		ProviderID:      &providerID,
		ExternalOrderID: "EXT-100",
	}

	orderRepo := &fakeOrderRepository{
		orders:         []*orderentity.Order{order},
		completeResult: false,
	}

	providerRepo := &fakeProviderRepository{
		provider: newProvider(7, "provider-a"),
	}

	provider := &fakeSMMProvider{
		status: orderentity.OrderStatusCompleted,
	}

	checkoutRepo := &fakeCheckoutTransactionRepository{}

	job, notificationRepo, _ := newStatusSyncJob(
		orderRepo,
		providerRepo,
		provider,
		checkoutRepo,
	)

	err := job.Run(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if order.Status != orderentity.OrderStatusProcessing {
		t.Fatalf(
			"expected order to remain PROCESSING, got %s",
			order.Status,
		)
	}

	if len(notificationRepo.created) != 0 {
		t.Fatalf(
			"expected no completion notification, got %d",
			len(notificationRepo.created),
		)
	}

	if len(checkoutRepo.calls) != 0 {
		t.Fatalf(
			"expected no refund, got %d",
			len(checkoutRepo.calls),
		)
	}
}
