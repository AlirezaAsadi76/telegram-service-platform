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

type fakePaymentRepository struct {
	paymentByIdempotency map[string]*paymententity.Payment
	createdPayment       *paymententity.Payment
	createErr            error
	getErr               error
}

func newFakePaymentRepository() *fakePaymentRepository {
	return &fakePaymentRepository{
		paymentByIdempotency: make(
			map[string]*paymententity.Payment,
		),
	}
}

func (f *fakePaymentRepository) Create(_ context.Context, payment *paymententity.Payment) error {
	if f.createErr != nil {
		return f.createErr
	}

	payment.ID = 100
	f.createdPayment = payment
	f.paymentByIdempotency[payment.IdempotencyKey] = payment
	return nil
}

func (f *fakePaymentRepository) GetByID(_ context.Context, _ uint64) (*paymententity.Payment, error) {
	return nil, nil
}

func (f *fakePaymentRepository) GetByOrderID(_ context.Context, _ uint64) (*paymententity.Payment, error) {
	return nil, nil
}

func (f *fakePaymentRepository) GetByIdempotencyKey(_ context.Context, key string) (*paymententity.Payment, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	payment, ok := f.paymentByIdempotency[key]
	if !ok {
		return nil, richerror.New(
			"fake.get_by_idempotency_key",
			nil,
		).WithKind(richerror.KindNotFound)
	}

	return payment, nil
}

func (f *fakePaymentRepository) UpdateStatus(_ context.Context, _ uint64, _ paymententity.PaymentStatus) error {
	return nil
}

func (f *fakePaymentRepository) GetPending(_ context.Context) ([]paymententity.Payment, error) {
	return nil, nil
}

func (f *fakePaymentRepository) GetExpired(_ context.Context) ([]paymententity.Payment, error) {
	return nil, nil
}

func TestService_CreateIntent(t *testing.T) {
	repo := newFakePaymentRepository()

	service := paymentservice.New(
		repo,
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
		t.Fatalf(
			"expected payment ID 100, got %d",
			resp.PaymentID,
		)
	}

	if resp.Status != paymententity.PaymentStatusCreating {
		t.Fatalf(
			"expected status CREATING, got %s",
			resp.Status,
		)
	}

	if repo.createdPayment == nil {
		t.Fatal("expected payment to be created")
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
