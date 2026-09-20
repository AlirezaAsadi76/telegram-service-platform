package orderfulfillertesting

import (
	"context"
	"errors"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
	"testing"
	"time"

	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/entity/productentity"
	"telegram-service-platform/params/smmparams"
	"telegram-service-platform/scheduler/jobs/orderfulfillerjob"
	"telegram-service-platform/service/notificationservice"
	"telegram-service-platform/service/orderservice"
	"telegram-service-platform/service/smmproviderservice"
)

func TestJob_Run_CreatedSMMOrder(t *testing.T) {
	// Arrange

	order := &orderentity.Order{
		ID:          42,
		UserID:      100,
		ProductType: productentity.ProductTypeSMM,
		ProductID:   55,
		Quantity:    1000,
		TargetLink:  "https://example.com/test",
		Status:      orderentity.OrderStatusPaid,
	}

	orderRepo := &fakeOrderRepository{
		order:       order,
		claimResult: true,
	}

	orderService := orderservice.New(orderRepo)

	providerRepo := &fakeProviderRepository{
		provider: newProvider(7, "provider-a"),
	}

	provider := &fakeSMMProvider{
		response: smmparams.CreateOrderAdapterResponse{
			Outcome:         smmparams.CreateOrderOutcomeCreated,
			ExternalOrderID: "EXT-100",
		},
	}

	smmService := smmproviderservice.New(
		providerRepo,
		smmproviderservice.Config{
			FailureThreshold: 3,
			SuccessThreshold: 1,
		},
	)

	smmService.RegisterProvider(
		"provider-a",
		provider,
	)

	redis := &fakeRedis{
		result: []string{
			"orders:paid",
			"42",
		},
	}

	notificationRepo := &fakeNotificationRepository{}
	notificationRedis := &fakeNotificationRedis{}

	notificationService := notificationservice.New(
		notificationRepo,
		notificationRedis,
		notificationservice.Config{},
	)

	job := orderfulfillerjob.New(
		orderService,
		smmService,
		notificationService,
		nil,
		nil,
		redis,
		orderfulfillerjob.Config{
			QueueKey: "orders:paid",
			Timeout:  time.Second,
		},
	)

	err := job.Run(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if order.Status != orderentity.OrderStatusProcessing {
		t.Fatalf(
			"expected order status PROCESSING, got %s",
			order.Status,
		)
	}

	if order.ProviderID == nil {
		t.Fatal("expected provider ID to be persisted")
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

	if provider.createCall != 1 {
		t.Fatalf(
			"expected provider Create to be called once, got %d",
			provider.createCall,
		)
	}

	if len(orderRepo.attempts) != 1 {
		t.Fatalf(
			"expected one fulfillment attempt, got %d",
			len(orderRepo.attempts),
		)
	}

	attempt := orderRepo.attempts[0]

	if attempt.Outcome != orderentity.FulfillmentAttemptOutcomeCreated {
		t.Fatalf(
			"expected attempt outcome CREATED, got %s",
			attempt.Outcome,
		)
	}

	if attempt.ExternalOrderID != "EXT-100" {
		t.Fatalf(
			"expected attempt external order ID EXT-100, got %s",
			attempt.ExternalOrderID,
		)
	}

	if attempt.ResolvedAt == nil {
		t.Fatal("expected fulfillment attempt to be resolved")
	}

	if len(notificationRepo.created) != 1 {
		t.Fatalf(
			"expected one notification, got %d",
			len(notificationRepo.created),
		)
	}

	notification := notificationRepo.created[0]

	if notification.UserID != order.UserID {
		t.Fatalf(
			"expected notification user ID %d, got %d",
			order.UserID,
			notification.UserID,
		)
	}

	expectedEvents := []string{
		"get_order",
		"claim",
		"save_external_order",
		"create_attempt",
		"resolve_attempt",
	}

	for i, expected := range expectedEvents {
		if i >= len(orderRepo.events) {
			t.Fatalf(
				"missing event %q",
				expected,
			)
		}

		if orderRepo.events[i] != expected {
			t.Fatalf(
				"expected event %q at position %d, got %q",
				expected,
				i,
				orderRepo.events[i],
			)
		}
	}
}

func TestJob_Run_SkipsAlreadyProcessedOrder(t *testing.T) {
	// Arrange

	order := &orderentity.Order{
		ID:          42,
		UserID:      100,
		ProductType: productentity.ProductTypeSMM,
		Status:      orderentity.OrderStatusProcessing,
	}

	orderRepo := &fakeOrderRepository{
		order: order,
	}

	orderService := orderservice.New(orderRepo)

	providerRepo := &fakeProviderRepository{
		provider: newProvider(7, "provider-a"),
	}

	provider := &fakeSMMProvider{
		response: smmparams.CreateOrderAdapterResponse{
			Outcome:         smmparams.CreateOrderOutcomeCreated,
			ExternalOrderID: "EXT-200",
		},
	}

	smmService := smmproviderservice.New(
		providerRepo,
		smmproviderservice.Config{
			FailureThreshold: 3,
			SuccessThreshold: 1,
		},
	)

	smmService.RegisterProvider(
		"provider-a",
		provider,
	)

	redis := &fakeRedis{
		result: []string{
			"orders:paid",
			"42",
		},
	}

	job := orderfulfillerjob.New(
		orderService,
		smmService,
		nil,
		nil,
		nil,
		redis,
		orderfulfillerjob.Config{
			QueueKey: "orders:paid",
			Timeout:  time.Second,
		},
	)

	// Act

	err := job.Run(context.Background())

	// Assert

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if provider.createCall != 0 {
		t.Fatalf(
			"expected provider Create not to be called, got %d",
			provider.createCall,
		)
	}

	if len(orderRepo.events) != 1 ||
		orderRepo.events[0] != "get_order" {
		t.Fatalf(
			"expected only GetByID, got %#v",
			orderRepo.events,
		)
	}
}

func TestJob_Run_UnknownProviderResultCreatesRecoveryAttempt(t *testing.T) {
	// Arrange

	order := &orderentity.Order{
		ID:          42,
		UserID:      100,
		ProductType: productentity.ProductTypeSMM,
		Status:      orderentity.OrderStatusPaid,
	}

	orderRepo := &fakeOrderRepository{
		order:       order,
		claimResult: true,
	}

	orderService := orderservice.New(orderRepo)

	providerRepo := &fakeProviderRepository{
		provider: newProvider(7, "provider-a"),
	}

	provider := &fakeSMMProvider{
		response: smmparams.CreateOrderAdapterResponse{
			Outcome: smmparams.CreateOrderOutcomeUnknown,
		},
	}

	smmService := smmproviderservice.New(
		providerRepo,
		smmproviderservice.Config{
			FailureThreshold: 3,
			SuccessThreshold: 1,
		},
	)

	smmService.RegisterProvider(
		"provider-a",
		provider,
	)

	redis := &fakeRedis{
		result: []string{
			"orders:paid",
			"42",
		},
	}

	job := orderfulfillerjob.New(
		orderService,
		smmService,
		nil,
		nil,
		nil,
		redis,
		orderfulfillerjob.Config{
			QueueKey: "orders:paid",
			Timeout:  time.Second,
		},
	)

	// Act

	err := job.Run(context.Background())

	// Assert

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if order.Status != orderentity.OrderStatusProcessing {
		t.Fatalf(
			"expected PROCESSING, got %s",
			order.Status,
		)
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

	if len(orderRepo.attempts) != 1 {
		t.Fatalf(
			"expected one fulfillment attempt, got %d",
			len(orderRepo.attempts),
		)
	}

	attempt := orderRepo.attempts[0]

	if attempt.Outcome != orderentity.FulfillmentAttemptOutcomeUnknown {
		t.Fatalf(
			"expected UNKNOWN attempt, got %s",
			attempt.Outcome,
		)
	}

	if provider.createCall != 1 {
		t.Fatalf(
			"expected provider Create once, got %d",
			provider.createCall,
		)
	}
}

func TestJob_Run_ClaimFailureDoesNotCallProvider(t *testing.T) {
	// Arrange

	order := &orderentity.Order{
		ID:          42,
		UserID:      100,
		ProductType: productentity.ProductTypeSMM,
		Status:      orderentity.OrderStatusPaid,
	}

	orderRepo := &fakeOrderRepository{
		order: order,
		claimErr: richerror.New(
			"orderservice.ClaimForProcessing",
			errors.New("fake claim error"),
		).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryFailed),
	}

	orderService := orderservice.New(orderRepo)

	providerRepo := &fakeProviderRepository{
		provider: newProvider(7, "provider-a"),
	}

	provider := &fakeSMMProvider{
		response: smmparams.CreateOrderAdapterResponse{
			Outcome:         smmparams.CreateOrderOutcomeCreated,
			ExternalOrderID: "EXT-200",
		},
	}

	smmService := smmproviderservice.New(
		providerRepo,
		smmproviderservice.Config{
			FailureThreshold: 3,
			SuccessThreshold: 1,
		},
	)

	smmService.RegisterProvider(
		"provider-a",
		provider,
	)

	redis := &fakeRedis{
		result: []string{
			"orders:paid",
			"42",
		},
	}

	job := orderfulfillerjob.New(
		orderService,
		smmService,
		nil,
		nil,
		nil,
		redis,
		orderfulfillerjob.Config{
			QueueKey: "orders:paid",
			Timeout:  time.Second,
		},
	)

	// Act

	err := job.Run(context.Background())

	// Assert

	if err == nil {
		t.Fatal("expected claim error, got nil")
	}

	if provider.createCall != 0 {
		t.Fatalf(
			"expected provider Create not to be called, got %d",
			provider.createCall,
		)
	}

	expectedEvents := []string{
		"get_order",
		"claim",
	}

	if len(orderRepo.events) != len(expectedEvents) {
		t.Fatalf(
			"expected %d events, got %d: %#v",
			len(expectedEvents),
			len(orderRepo.events),
			orderRepo.events,
		)
	}

	for i, expected := range expectedEvents {
		if orderRepo.events[i] != expected {
			t.Fatalf(
				"expected event %q at position %d, got %q",
				expected,
				i,
				orderRepo.events[i],
			)
		}
	}
}
