package paymenttesting

import (
	"context"
	"errors"
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/params/paymentparams"
	"telegram-service-platform/params/paymentproviderparams"
	"telegram-service-platform/pkg/richerror"
	"telegram-service-platform/service/paymentservice"
	"testing"

	"github.com/shopspring/decimal"
)

func TestService_Initiate_Success(t *testing.T) {
	repo := newFakePaymentRepository()

	payment := &paymententity.Payment{
		ID:      100,
		OrderID: 10,
		UserID:  20,
		Method:  paymententity.PaymentMethodZarinpal,
		Status:  paymententity.PaymentStatusCreating,
	}

	repo.payments[payment.ID] = payment

	provider := &fakePaymentProvider{
		createResponse: paymentproviderparams.CreateResponse{
			ExternalID: "EXT-123",
			PaymentURL: "https://provider/pay/123",
		},
	}

	service := paymentservice.New(
		repo,
		newFakePaymentConfirmationRepository(),
		provider,
		nil,
	)

	resp, err := service.Initiate(
		context.Background(),
		paymentparams.InitiateRequest{
			PaymentID:   100,
			CallbackURL: "https://example.com/payment/callback",
			Description: "Order #10",
		},
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Status != paymententity.PaymentStatusPending {
		t.Fatalf(
			"expected PENDING, got %s",
			resp.Status,
		)
	}

	if resp.ExternalID != "EXT-123" {
		t.Fatalf("unexpected external ID")
	}

	if resp.PaymentURL != "https://provider/pay/123" {
		t.Fatalf("unexpected payment URL")
	}

	if provider.createCalls != 1 {
		t.Fatalf(
			"expected 1 provider create call, got %d",
			provider.createCalls,
		)
	}

	if repo.markInitiatedCalls != 1 {
		t.Fatalf(
			"expected 1 mark initiated call, got %d",
			repo.markInitiatedCalls,
		)
	}
}

func TestService_Initiate_ProviderTimeout(t *testing.T) {
	repo := newFakePaymentRepository()

	payment := &paymententity.Payment{
		ID:       100,
		OrderID:  10,
		Status:   paymententity.PaymentStatusCreating,
		Method:   paymententity.PaymentMethodZarinpal,
		Amount:   entity.Amount(decimal.NewFromInt(100000)),
		Currency: entity.CurrencyTOMAN,
	}

	repo.payments[payment.ID] = payment

	provider := &fakePaymentProvider{
		createErr: richerror.New(
			"fakeprovider.create",
			context.DeadlineExceeded,
		).WithCode(richerror.CodePaymentProviderTimeout),
	}

	service := paymentservice.New(
		repo,
		newFakePaymentConfirmationRepository(),
		provider,
		nil,
	)

	_, err := service.Initiate(
		context.Background(),
		paymentparams.InitiateRequest{
			PaymentID: payment.ID,
		},
	)

	if err == nil {
		t.Fatal("expected error")
	}

	persisted := repo.payments[payment.ID]

	if persisted.Status != paymententity.PaymentStatusUnknown {
		t.Fatalf("expected UNKNOWN, got %s", persisted.Status)
	}
}

func TestService_Initiate_MarkInitiatedFailure(t *testing.T) {

	repo := newFakePaymentRepository()

	payment := &paymententity.Payment{
		ID:       100,
		OrderID:  10,
		UserID:   20,
		Method:   paymententity.PaymentMethodZarinpal,
		Amount:   entity.Amount(decimal.NewFromInt(100000)),
		Currency: entity.CurrencyTOMAN,
		Status:   paymententity.PaymentStatusCreating,
	}

	repo.payments[payment.ID] = payment

	markInitiatedErr := errors.New(
		"database update failed",
	)

	repo.markInitiatedErr = markInitiatedErr

	provider := &fakePaymentProvider{
		createResponse: paymentproviderparams.CreateResponse{
			ExternalID: "EXT-123",
			PaymentURL: "https://provider.test/pay/123",
		},
	}

	service := paymentservice.New(
		repo,
		newFakePaymentConfirmationRepository(),
		provider,
		nil,
	)

	response, err := service.Initiate(
		context.Background(),
		paymentparams.InitiateRequest{
			PaymentID:   payment.ID,
			CallbackURL: "https://example.com/payment/callback",
			Description: "Order #10",
		},
	)

	if err == nil {
		t.Fatal("expected Initiate to return an error")
	}

	if response != nil {
		t.Fatal(
			"expected response to be nil when MarkInitiated fails",
		)
	}

	if !richerror.IsCode(err, richerror.CodePaymentInitiationUpdateFailed) {
		t.Fatalf(
			"expected error code %s",
			richerror.CodePaymentInitiationUpdateFailed,
		)
	}

	if provider.createCalls != 1 {
		t.Fatalf(
			"expected provider Create to be called once, got %d",
			provider.createCalls,
		)
	}

	if repo.markInitiatedCalls != 1 {
		t.Fatalf(
			"expected MarkInitiated to be called once, got %d",
			repo.markInitiatedCalls,
		)
	}

	if repo.updateStatusCalls != 0 {
		t.Fatalf(
			"expected UpdateStatus not to be called, got %d calls",
			repo.updateStatusCalls,
		)
	}

	persisted := repo.payments[payment.ID]

	if persisted == nil {
		t.Fatal("expected payment to remain in repository")
	}

	if persisted.Status != paymententity.PaymentStatusCreating {
		t.Fatalf(
			"expected payment to remain CREATING, got %s",
			persisted.Status,
		)
	}
}

func TestService_Initiate_InvalidProviderResponse(t *testing.T) {

	repo := newFakePaymentRepository()

	payment := &paymententity.Payment{
		ID:       100,
		OrderID:  10,
		UserID:   20,
		Method:   paymententity.PaymentMethodZarinpal,
		Amount:   entity.Amount(decimal.NewFromInt(100000)),
		Currency: entity.CurrencyTOMAN,
		Status:   paymententity.PaymentStatusCreating,
	}

	repo.payments[payment.ID] = payment

	provider := &fakePaymentProvider{
		createResponse: paymentproviderparams.CreateResponse{
			ExternalID: "",
			PaymentURL: "https://provider.test/pay/100",
		},
	}

	service := paymentservice.New(
		repo,
		newFakePaymentConfirmationRepository(),
		provider,
		nil,
	)

	req := paymentparams.InitiateRequest{
		PaymentID:   payment.ID,
		CallbackURL: "https://example.com/payments/zarinpal/callback",
		Description: "Invalid provider response test",
	}

	response, err := service.Initiate(
		context.Background(),
		req,
	)

	if err == nil {
		t.Fatal("expected Initiate to return an error")
	}

	if response != nil {
		t.Fatal(
			"expected response to be nil for invalid provider response",
		)
	}

	if !richerror.IsCode(err, richerror.CodePaymentProviderInvalidResponse) {
		t.Fatalf(
			"expected error code %s",
			richerror.CodePaymentProviderInvalidResponse,
		)
	}

	if provider.createCalls != 1 {
		t.Fatalf(
			"expected provider Create to be called once, got %d",
			provider.createCalls,
		)
	}

	if repo.markInitiatedCalls != 0 {
		t.Fatalf(
			"expected MarkInitiated not to be called, got %d",
			repo.markInitiatedCalls,
		)
	}

	persisted := repo.payments[payment.ID]

	if persisted == nil {
		t.Fatal("expected payment to remain in repository")
	}

	if persisted.Status != paymententity.PaymentStatusUnknown {
		t.Fatalf(
			"expected payment status UNKNOWN, got %s",
			persisted.Status,
		)
	}

	if persisted.ExternalID != "" {
		t.Fatalf(
			"expected external ID to remain empty, got %s",
			persisted.ExternalID,
		)

	}

	if persisted.PaymentURL != "" {
		t.Fatalf(
			"expected payment URL to remain empty, got %s",
			persisted.PaymentURL,
		)
	}
}
