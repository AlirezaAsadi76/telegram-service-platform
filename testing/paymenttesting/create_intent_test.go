package paymenttesting

import (
	"context"
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/params/paymentparams"
	"telegram-service-platform/pkg/richerror"
	"telegram-service-platform/service/paymentservice"
	"testing"

	"github.com/shopspring/decimal"
)

func TestService_CreateIntent_Success(t *testing.T) {
	repo := newFakePaymentRepository()
	confirmationRepo := newFakePaymentConfirmationRepository()

	service := paymentservice.New(
		repo,
		confirmationRepo,
		nil,
		nil,
	)

	req := paymentparams.CreateIntentRequest{
		OrderID:        10,
		UserID:         20,
		Method:         paymententity.PaymentMethodZarinpal,
		Amount:         entity.Amount(decimal.NewFromInt(50000)),
		Currency:       entity.CurrencyTOMAN,
		IdempotencyKey: "checkout:order:10:attempt:abc",
	}

	resp, err := service.CreateIntent(
		context.Background(),
		req,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp == nil {
		t.Fatal("expected response")
	}

	if resp.PaymentID != 100 {
		t.Fatalf("expected payment ID 100, got %d", resp.PaymentID)
	}

	if resp.Status != paymententity.PaymentStatusCreating {
		t.Fatalf(
			"expected CREATING, got %s",
			resp.Status,
		)
	}

	if repo.createCalls != 1 {
		t.Fatalf(
			"expected 1 create call, got %d",
			repo.createCalls,
		)
	}

	if repo.createdPayment.OrderID != req.OrderID {
		t.Fatalf("unexpected order ID")
	}

	if repo.createdPayment.UserID != req.UserID {
		t.Fatalf("unexpected user ID")
	}

	if repo.createdPayment.IdempotencyKey != req.IdempotencyKey {
		t.Fatalf("unexpected idempotency key")
	}
}

func TestService_CreateIntent_RecoverExistingPaymentAfterConflict(
	t *testing.T,
) {
	repo := newFakePaymentRepository()
	confirmationRepo := newFakePaymentConfirmationRepository()
	existing := &paymententity.Payment{
		ID:             100,
		OrderID:        10,
		UserID:         20,
		Method:         paymententity.PaymentMethodZarinpal,
		Amount:         entity.Amount(decimal.NewFromInt(50000)),
		Currency:       entity.CurrencyTOMAN,
		Status:         paymententity.PaymentStatusCreating,
		IdempotencyKey: "checkout:order:10:attempt:abc",
	}

	repo.paymentByIdempotency[existing.IdempotencyKey] = existing

	repo.getByIdempotencySequence = []error{
		richerror.New(
			"fake.get_by_idempotency_key",
			nil,
		).WithKind(richerror.KindNotFound),

		nil,
	}

	repo.createErr = richerror.New(
		"fake.create",
		nil,
	).
		WithKind(richerror.KindConflict).
		WithCode(
			richerror.CodePaymentIdempotencyKeyReused,
		)

	service := paymentservice.New(repo, confirmationRepo, nil, nil)

	req := paymentparams.CreateIntentRequest{
		OrderID:        10,
		UserID:         20,
		Method:         paymententity.PaymentMethodZarinpal,
		Amount:         entity.Amount(decimal.NewFromInt(50000)),
		Currency:       entity.CurrencyTOMAN,
		IdempotencyKey: existing.IdempotencyKey,
	}

	resp, err := service.CreateIntent(
		context.Background(),
		req,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp == nil {
		t.Fatal("expected response")
	}

	if resp.PaymentID != existing.ID {
		t.Fatalf(
			"expected payment ID %d, got %d",
			existing.ID,
			resp.PaymentID,
		)
	}

	if resp.Status != existing.Status {
		t.Fatalf(
			"expected status %s, got %s",
			existing.Status,
			resp.Status,
		)
	}

	if repo.getByIdempotencyCalls != 2 {
		t.Fatalf(
			"expected 2 idempotency lookups, got %d",
			repo.getByIdempotencyCalls,
		)
	}
}

func TestService_CreateIntent_ActivePaymentAlreadyExists(
	t *testing.T,
) {
	repo := newFakePaymentRepository()
	confirmationRepo := newFakePaymentConfirmationRepository()
	repo.getByIdempotencySequence = []error{
		richerror.New(
			"fake.get_by_idempotency_key",
			nil,
		).WithKind(richerror.KindNotFound),
	}

	repo.createErr = richerror.New(
		"fake.create",
		nil,
	).
		WithKind(richerror.KindConflict).
		WithCode(
			richerror.CodePaymentIntentAlreadyExists,
		)

	service := paymentservice.New(repo, confirmationRepo, nil, nil)

	req := paymentparams.CreateIntentRequest{
		OrderID:        10,
		UserID:         20,
		Method:         paymententity.PaymentMethodZarinpal,
		Amount:         entity.Amount(decimal.NewFromInt(50000)),
		Currency:       entity.CurrencyTOMAN,
		IdempotencyKey: "another-key",
	}

	_, err := service.CreateIntent(
		context.Background(),
		req,
	)

	if err == nil {
		t.Fatal("expected error")
	}

	if !richerror.IsCode(
		err,
		richerror.CodePaymentIntentAlreadyExists,
	) {
		t.Fatalf(
			"expected PAYMENT_INTENT_ALREADY_EXISTS, got %v",
			err,
		)
	}
}
