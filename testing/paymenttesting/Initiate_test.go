package paymenttesting

import (
	"context"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/params/paymentparams"
	"telegram-service-platform/params/paymentproviderparams"
	"telegram-service-platform/service/paymentservice"
	"testing"
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
