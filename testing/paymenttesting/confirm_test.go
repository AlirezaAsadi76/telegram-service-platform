package paymenttesting

import (
	"context"
	"testing"

	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/params/paymentparams"
	"telegram-service-platform/params/paymentproviderparams"
	"telegram-service-platform/service/paymentservice"
)

func TestService_ConfirmPayment_Success(t *testing.T) {
	repo := newFakePaymentRepository()
	confirmationRepo := newFakePaymentConfirmationRepository()

	payment := &paymententity.Payment{
		ID:         100,
		OrderID:    10,
		UserID:     20,
		Method:     paymententity.PaymentMethodZarinpal,
		Status:     paymententity.PaymentStatusPending,
		ExternalID: "EXT-100",
	}

	repo.payments[payment.ID] = payment

	provider := &fakePaymentProvider{
		verifyResponse: paymentproviderparams.VerifyResponse{
			Status: paymententity.PaymentStatusSuccess,
		},
	}

	service := paymentservice.New(
		repo,
		confirmationRepo,
		provider,
		nil,
	)

	resp, err := service.ConfirmPayment(
		context.Background(),
		paymentparams.ConfirmPaymentRequest{
			PaymentID:    payment.ID,
			CallbackData: map[string]any{"authority": "ABC"},
		},
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.PaymentID != payment.ID {
		t.Fatalf("unexpected payment ID")
	}

	if resp.OrderID != payment.OrderID {
		t.Fatalf("unexpected order ID")
	}

	if resp.Status != paymententity.PaymentStatusSuccess {
		t.Fatalf(
			"expected SUCCESS, got %s",
			resp.Status,
		)
	}

	if provider.verifyCalls != 1 {
		t.Fatalf(
			"expected 1 verify call, got %d",
			provider.verifyCalls,
		)
	}

	if confirmationRepo.confirmCalls != 1 {
		t.Fatalf(
			"expected 1 confirmation call, got %d",
			confirmationRepo.confirmCalls,
		)
	}
}

func TestService_ConfirmPayment_ProviderFailed(t *testing.T) {
	repo := newFakePaymentRepository()
	confirmationRepo := newFakePaymentConfirmationRepository()

	payment := &paymententity.Payment{
		ID:      100,
		OrderID: 10,
		Status:  paymententity.PaymentStatusPending,
	}

	repo.payments[payment.ID] = payment

	provider := &fakePaymentProvider{
		verifyResponse: paymentproviderparams.VerifyResponse{
			Status: paymententity.PaymentStatusFailed,
		},
	}

	service := paymentservice.New(
		repo,
		confirmationRepo,
		provider,
		nil,
	)

	resp, err := service.ConfirmPayment(
		context.Background(),
		paymentparams.ConfirmPaymentRequest{
			PaymentID: payment.ID,
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Status != paymententity.PaymentStatusFailed {
		t.Fatalf(
			"expected status %s, got %s",
			paymententity.PaymentStatusFailed,
			resp.Status,
		)
	}

	if confirmationRepo.confirmCalls != 0 {
		t.Fatalf(
			"expected confirm not to be called, got %d",
			confirmationRepo.confirmCalls,
		)
	}

	if confirmationRepo.failCalls != 1 {
		t.Fatalf(
			"expected fail to be called once, got %d",
			confirmationRepo.failCalls,
		)
	}
}

func TestService_ConfirmPayment_AlreadySuccess(t *testing.T) {
	repo := newFakePaymentRepository()
	confirmationRepo := newFakePaymentConfirmationRepository()

	payment := &paymententity.Payment{
		ID:      100,
		OrderID: 10,
		Status:  paymententity.PaymentStatusSuccess,
	}

	repo.payments[payment.ID] = payment

	provider := &fakePaymentProvider{}

	service := paymentservice.New(
		repo,
		confirmationRepo,
		provider,
		nil,
	)

	resp, err := service.ConfirmPayment(
		context.Background(),
		paymentparams.ConfirmPaymentRequest{
			PaymentID: payment.ID,
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Status != paymententity.PaymentStatusSuccess {
		t.Fatalf("expected success status")
	}

	if provider.verifyCalls != 0 {
		t.Fatalf(
			"expected provider not to be called, got %d calls",
			provider.verifyCalls,
		)
	}

	if confirmationRepo.confirmCalls != 0 {
		t.Fatalf(
			"expected confirmation repository not to be called, got %d calls",
			confirmationRepo.confirmCalls,
		)
	}
}

func TestService_ConfirmPayment_InvalidState(t *testing.T) {
	repo := newFakePaymentRepository()
	confirmationRepo := newFakePaymentConfirmationRepository()

	payment := &paymententity.Payment{
		ID:      100,
		OrderID: 10,
		Status:  paymententity.PaymentStatusCanceled,
	}

	repo.payments[payment.ID] = payment

	provider := &fakePaymentProvider{}

	service := paymentservice.New(
		repo,
		confirmationRepo,
		provider,
		nil,
	)

	_, err := service.ConfirmPayment(
		context.Background(),
		paymentparams.ConfirmPaymentRequest{
			PaymentID: payment.ID,
		},
	)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if provider.verifyCalls != 0 {
		t.Fatalf(
			"expected provider not to be called, got %d calls",
			provider.verifyCalls,
		)
	}

	if confirmationRepo.confirmCalls != 0 {
		t.Fatalf(
			"expected confirmation repository not to be called, got %d calls",
			confirmationRepo.confirmCalls,
		)
	}
}

func TestService_ConfirmPayment_VerificationError(t *testing.T) {
	repo := newFakePaymentRepository()
	confirmationRepo := newFakePaymentConfirmationRepository()

	payment := &paymententity.Payment{
		ID:      100,
		OrderID: 10,
		Status:  paymententity.PaymentStatusPending,
	}

	repo.payments[payment.ID] = payment

	provider := &fakePaymentProvider{
		verifyErr: context.DeadlineExceeded,
	}

	service := paymentservice.New(
		repo,
		confirmationRepo,
		provider,
		nil,
	)

	_, err := service.ConfirmPayment(
		context.Background(),
		paymentparams.ConfirmPaymentRequest{
			PaymentID: payment.ID,
		},
	)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if confirmationRepo.confirmCalls != 0 {
		t.Fatalf(
			"expected confirmation repository not to be called, got %d",
			confirmationRepo.confirmCalls,
		)
	}
}
