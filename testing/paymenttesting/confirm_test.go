package paymenttesting

import (
	"context"
	"testing"

	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/params/paymentparams"
	"telegram-service-platform/params/paymentproviderparams"
	"telegram-service-platform/service/paymentservice"
)

type fakePaymentProvider1 struct {
	verifyResponse paymentproviderparams.VerifyResponse
	verifyErr      error
	verifyCalls    int
}

func (f *fakePaymentProvider1) Create(
	_ context.Context,
	_ paymentproviderparams.CreateRequest,
) (paymentproviderparams.CreateResponse, error) {
	return paymentproviderparams.CreateResponse{}, nil
}

func (f *fakePaymentProvider1) Verify(
	_ context.Context,
	_ paymentproviderparams.VerifyRequest,
) (paymentproviderparams.VerifyResponse, error) {
	f.verifyCalls++

	return f.verifyResponse, f.verifyErr
}

func TestService_ConfirmPayment_Success(t *testing.T) {
	repo := newFakePaymentRepository()
	confirmationRepo := newFakePaymentConfirmationRepository()

	payment := &paymententity.Payment{
		ID:      100,
		OrderID: 10,
		Status:  paymententity.PaymentStatusPending,
	}

	repo.payments[payment.ID] = payment

	provider := &fakePaymentProvider1{
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
			ExternalID:   "external-100",
			CallbackData: map[string]any{"authority": "abc"},
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.PaymentID != payment.ID {
		t.Fatalf("expected payment ID %d, got %d", payment.ID, resp.PaymentID)
	}

	if resp.OrderID != payment.OrderID {
		t.Fatalf("expected order ID %d, got %d", payment.OrderID, resp.OrderID)
	}

	if resp.Status != paymententity.PaymentStatusSuccess {
		t.Fatalf(
			"expected status %s, got %s",
			paymententity.PaymentStatusSuccess,
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
			"expected 1 confirm call, got %d",
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

	provider := &fakePaymentProvider1{
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
			"expected confirmation repository not to be called, got %d",
			confirmationRepo.confirmCalls,
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

	provider := &fakePaymentProvider1{}

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

	provider := &fakePaymentProvider1{}

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

	provider := &fakePaymentProvider1{
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
