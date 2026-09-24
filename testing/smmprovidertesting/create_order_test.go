package smmprovidertesting

import (
	"context"
	"testing"

	"telegram-service-platform/entity/providerentity"
	"telegram-service-platform/params/smmparams"
	"telegram-service-platform/pkg/richerror"
)

func TestService_CreateOrder_Created(t *testing.T) {
	repo := &fakeProviderRepository{
		providers: []*providerentity.Provider{
			newSMMProvider(1, "provider-a"),
		},
	}

	provider := &fakeSMMProvider{
		createResponse: smmparams.CreateOrderAdapterResponse{
			Outcome:         smmparams.CreateOrderOutcomeCreated,
			ExternalOrderID: "EXT-100",
		},
	}

	service := newTestService(
		repo,
		map[string]*fakeSMMProvider{
			"provider-a": provider,
		},
	)

	response, err := service.CreateOrder(
		context.Background(),
		smmparams.CreateOrderAdapterRequest{},
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response.Outcome != smmparams.CreateOrderOutcomeCreated {
		t.Fatalf(
			"expected CREATED, got %s",
			response.Outcome,
		)
	}

	if response.ProviderID != 1 {
		t.Fatalf(
			"expected provider ID 1, got %d",
			response.ProviderID,
		)
	}

	if response.ProviderName != "provider-a" {
		t.Fatalf(
			"expected provider name provider-a, got %s",
			response.ProviderName,
		)
	}

	if response.ExternalOrderID != "EXT-100" {
		t.Fatalf(
			"expected external order ID EXT-100, got %s",
			response.ExternalOrderID,
		)
	}

	if provider.createCalls != 1 {
		t.Fatalf(
			"expected Create to be called once, got %d",
			provider.createCalls,
		)
	}
}

func TestService_CreateOrder_RejectedThenCreated(t *testing.T) {
	repo := &fakeProviderRepository{
		providers: []*providerentity.Provider{
			newSMMProvider(1, "provider-a"),
			newSMMProvider(2, "provider-b"),
		},
	}

	firstProvider := &fakeSMMProvider{
		createResponse: smmparams.CreateOrderAdapterResponse{
			Outcome: smmparams.CreateOrderOutcomeRejected,
		},
	}

	secondProvider := &fakeSMMProvider{
		createResponse: smmparams.CreateOrderAdapterResponse{
			Outcome:         smmparams.CreateOrderOutcomeCreated,
			ExternalOrderID: "EXT-200",
		},
	}

	service := newTestService(
		repo,
		map[string]*fakeSMMProvider{
			"provider-a": firstProvider,
			"provider-b": secondProvider,
		},
	)

	response, err := service.CreateOrder(
		context.Background(),
		smmparams.CreateOrderAdapterRequest{},
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response.Outcome != smmparams.CreateOrderOutcomeCreated {
		t.Fatalf(
			"expected CREATED, got %s",
			response.Outcome,
		)
	}

	if response.ProviderID != 2 {
		t.Fatalf(
			"expected provider ID 2, got %d",
			response.ProviderID,
		)
	}

	if response.ExternalOrderID != "EXT-200" {
		t.Fatalf(
			"expected external order ID EXT-200, got %s",
			response.ExternalOrderID,
		)
	}

	if firstProvider.createCalls != 1 {
		t.Fatalf(
			"expected first provider to be called once, got %d",
			firstProvider.createCalls,
		)
	}

	if secondProvider.createCalls != 1 {
		t.Fatalf(
			"expected second provider to be called once, got %d",
			secondProvider.createCalls,
		)
	}
}

func TestService_CreateOrder_UnknownStopsFallback(t *testing.T) {
	repo := &fakeProviderRepository{
		providers: []*providerentity.Provider{
			newSMMProvider(1, "provider-a"),
			newSMMProvider(2, "provider-b"),
		},
	}

	firstProvider := &fakeSMMProvider{
		createResponse: smmparams.CreateOrderAdapterResponse{
			Outcome: smmparams.CreateOrderOutcomeUnknown,
		},
		createErr: richerror.New(
			"fakeprovider.Create",
			context.DeadlineExceeded,
		),
	}

	secondProvider := &fakeSMMProvider{
		createResponse: smmparams.CreateOrderAdapterResponse{
			Outcome:         smmparams.CreateOrderOutcomeCreated,
			ExternalOrderID: "EXT-300",
		},
	}

	service := newTestService(
		repo,
		map[string]*fakeSMMProvider{
			"provider-a": firstProvider,
			"provider-b": secondProvider,
		},
	)

	response, err := service.CreateOrder(
		context.Background(),
		smmparams.CreateOrderAdapterRequest{},
	)

	if err == nil {
		t.Fatal("expected error")
	}

	if response.Outcome != smmparams.CreateOrderOutcomeUnknown {
		t.Fatalf(
			"expected UNKNOWN, got %s",
			response.Outcome,
		)
	}

	if response.ProviderID != 1 {
		t.Fatalf(
			"expected provider ID 1, got %d",
			response.ProviderID,
		)
	}

	if firstProvider.createCalls != 1 {
		t.Fatalf(
			"expected first provider to be called once, got %d",
			firstProvider.createCalls,
		)
	}

	if secondProvider.createCalls != 0 {
		t.Fatalf(
			"expected second provider not to be called, got %d calls",
			secondProvider.createCalls,
		)
	}
}

func TestService_CreateOrder_CreatedWithoutExternalID(t *testing.T) {
	repo := &fakeProviderRepository{
		providers: []*providerentity.Provider{
			newSMMProvider(1, "provider-a"),
		},
	}

	provider := &fakeSMMProvider{
		createResponse: smmparams.CreateOrderAdapterResponse{
			Outcome: smmparams.CreateOrderOutcomeCreated,
		},
	}

	service := newTestService(
		repo,
		map[string]*fakeSMMProvider{
			"provider-a": provider,
		},
	)

	response, err := service.CreateOrder(
		context.Background(),
		smmparams.CreateOrderAdapterRequest{},
	)

	if err == nil {
		t.Fatal("expected error")
	}

	if response.Outcome != smmparams.CreateOrderOutcomeUnknown {
		t.Fatalf(
			"expected UNKNOWN, got %s",
			response.Outcome,
		)
	}

	if response.ProviderID != 1 {
		t.Fatalf(
			"expected provider ID 1, got %d",
			response.ProviderID,
		)
	}

	if !richerror.IsCode(
		err,
		richerror.CodeSMMProviderInvalidResponse,
	) {
		t.Fatalf(
			"expected error code %s",
			richerror.CodeSMMProviderInvalidResponse,
		)
	}

	if provider.createCalls != 1 {
		t.Fatalf(
			"expected provider to be called once, got %d",
			provider.createCalls,
		)
	}
}

func TestService_CreateOrder_AllProvidersRejected(t *testing.T) {
	repo := &fakeProviderRepository{
		providers: []*providerentity.Provider{
			newSMMProvider(1, "provider-a"),
			newSMMProvider(2, "provider-b"),
		},
	}

	firstProvider := &fakeSMMProvider{
		createResponse: smmparams.CreateOrderAdapterResponse{
			Outcome: smmparams.CreateOrderOutcomeRejected,
		},
	}

	secondProvider := &fakeSMMProvider{
		createResponse: smmparams.CreateOrderAdapterResponse{
			Outcome: smmparams.CreateOrderOutcomeRejected,
		},
	}

	service := newTestService(
		repo,
		map[string]*fakeSMMProvider{
			"provider-a": firstProvider,
			"provider-b": secondProvider,
		},
	)

	response, err := service.CreateOrder(
		context.Background(),
		smmparams.CreateOrderAdapterRequest{},
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response.Outcome != smmparams.CreateOrderOutcomeRejected {
		t.Fatalf(
			"expected REJECTED, got %s",
			response.Outcome,
		)
	}

	if firstProvider.createCalls != 1 {
		t.Fatalf(
			"expected first provider to be called once, got %d",
			firstProvider.createCalls,
		)
	}

	if secondProvider.createCalls != 1 {
		t.Fatalf(
			"expected second provider to be called once, got %d",
			secondProvider.createCalls,
		)
	}
}

func TestService_CreateOrder_TargetsRequestedProvider(
	t *testing.T,
) {
	firstProvider := &fakeSMMProvider{
		createResponse: smmparams.CreateOrderAdapterResponse{
			Outcome:         smmparams.CreateOrderOutcomeCreated,
			ExternalOrderID: "WRONG-100",
		},
	}

	japProvider := &fakeSMMProvider{
		createResponse: smmparams.CreateOrderAdapterResponse{
			Outcome:         smmparams.CreateOrderOutcomeCreated,
			ExternalOrderID: "JAP-100",
		},
	}

	repo := &fakeProviderRepository{
		providers: []*providerentity.Provider{
			newSMMProvider(1, "provider-a"),
			newSMMProvider(2, "justanotherpanel"),
		},
	}

	service := newTestService(
		repo,
		map[string]*fakeSMMProvider{
			"provider-a":       firstProvider,
			"justanotherpanel": japProvider,
		},
	)

	result, err := service.CreateOrder(
		context.Background(),
		smmparams.CreateOrderAdapterRequest{
			ProviderName: "justanotherpanel",
			ServiceID:    "123456",
			Link:         "https://example.com",
			Quantity:     1000,
		},
	)
	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if result.Outcome !=
		smmparams.CreateOrderOutcomeCreated {
		t.Fatalf(
			"expected CREATED, got %s",
			result.Outcome,
		)
	}

	if result.ProviderName != "justanotherpanel" {
		t.Fatalf(
			"expected justanotherpanel, got %s",
			result.ProviderName,
		)
	}

	if firstProvider.createCalls != 0 {
		t.Fatalf(
			"expected provider-a not to be called, got %d calls",
			firstProvider.createCalls,
		)
	}

	if japProvider.createCalls != 1 {
		t.Fatalf(
			"expected JAP provider to be called once, got %d calls",
			japProvider.createCalls,
		)
	}

	if japProvider.lastRequest.ServiceID != "123456" {
		t.Fatalf(
			"expected service ID 123456, got %s",
			japProvider.lastRequest.ServiceID,
		)
	}
}

func TestService_CreateOrder_TargetedProviderRejectedDoesNotFallback(
	t *testing.T,
) {
	targetProvider := &fakeSMMProvider{
		createResponse: smmparams.CreateOrderAdapterResponse{
			Outcome: smmparams.CreateOrderOutcomeRejected,
		},
	}

	fallbackProvider := &fakeSMMProvider{
		createResponse: smmparams.CreateOrderAdapterResponse{
			Outcome:         smmparams.CreateOrderOutcomeCreated,
			ExternalOrderID: "FALLBACK-100",
		},
	}

	repo := &fakeProviderRepository{
		providers: []*providerentity.Provider{
			newSMMProvider(1, "provider-a"),
			newSMMProvider(2, "provider-b"),
		},
	}

	service := newTestService(
		repo,
		map[string]*fakeSMMProvider{
			"provider-a": targetProvider,
			"provider-b": fallbackProvider,
		},
	)

	result, err := service.CreateOrder(
		context.Background(),
		smmparams.CreateOrderAdapterRequest{
			ProviderName: "provider-a",
			ServiceID:    "123456",
			Link:         "https://example.com",
			Quantity:     1000,
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.Outcome != smmparams.CreateOrderOutcomeRejected {
		t.Fatalf(
			"expected REJECTED, got %s",
			result.Outcome,
		)
	}

	if targetProvider.createCalls != 1 {
		t.Fatalf(
			"expected targeted provider to be called once, got %d",
			targetProvider.createCalls,
		)
	}

	if fallbackProvider.createCalls != 0 {
		t.Fatalf(
			"expected fallback provider not to be called, got %d",
			fallbackProvider.createCalls,
		)
	}

	if targetProvider.lastRequest.ProviderName != "provider-a" {
		t.Fatalf(
			"expected provider name provider-a, got %s",
			targetProvider.lastRequest.ProviderName,
		)
	}

	if targetProvider.lastRequest.ServiceID != "123456" {
		t.Fatalf(
			"expected service ID 123456, got %s",
			targetProvider.lastRequest.ServiceID,
		)
	}
}
