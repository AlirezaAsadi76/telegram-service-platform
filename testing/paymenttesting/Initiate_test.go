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
	provider := &fakePaymentProvider{
		response: paymentproviderparams.CreateResponse{
			ExternalID: "EXT-123",
			PaymentURL: "https://provider/pay/123",
		},
	}

	repo := newFakePaymentRepository()
	confirmationRepo := newFakePaymentConfirmationRepository()
	payment := &paymententity.Payment{
		ID:      100,
		OrderID: 10,
		UserID:  20,
		Method:  paymententity.PaymentMethodZarinpal,
		Status:  paymententity.PaymentStatusCreating,
	}

	repo.payments[100] = payment

	service := paymentservice.New(
		repo,
		confirmationRepo,
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
}
