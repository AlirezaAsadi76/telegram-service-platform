package paymenttesting

import (
	"context"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/params/paymentproviderparams"
	"telegram-service-platform/pkg/richerror"
)

type fakePaymentRepository struct {
	paymentByIdempotency map[string]*paymententity.Payment
	createdPayment       *paymententity.Payment
	payments             map[uint64]*paymententity.Payment
	createErr            error

	getSequence []error

	getByIdempotencyCalls int
	markInitiatedCalls    int
}

func newFakePaymentRepository() *fakePaymentRepository {
	return &fakePaymentRepository{
		paymentByIdempotency: make(
			map[string]*paymententity.Payment,
		),
		payments: make(map[uint64]*paymententity.Payment),
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
	payment, ok := f.payments[id]
	if !ok {
		return nil, richerror.New(
			"fake.get_by_id",
			nil,
		).
			WithKind(richerror.KindNotFound).
			WithCode(richerror.CodePaymentNotFound)
	}

	return payment, nil
}

func (f *fakePaymentRepository) GetByOrderID(_ context.Context, id uint64) (*paymententity.Payment, error) {
	return f.createdPayment, nil
}

func (f *fakePaymentRepository) GetByIdempotencyKey(_ context.Context, key string) (*paymententity.Payment, error) {
	index := f.getByIdempotencyCalls
	f.getByIdempotencyCalls++

	if index < len(f.getSequence) && f.getSequence[index] != nil {
		return nil, f.getSequence[index]
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

func (f *fakePaymentRepository) MarkInitiated(
	_ context.Context,
	paymentID uint64,
	status paymententity.PaymentStatus,
	externalID string,
	paymentURL string,
) error {
	payment, ok := f.payments[paymentID]
	if !ok {
		return richerror.New(
			"fake.mark_initiated",
			nil,
		).
			WithKind(richerror.KindNotFound).
			WithCode(richerror.CodePaymentNotFound)
	}

	if payment.Status != paymententity.PaymentStatusCreating {
		return richerror.New(
			"fake.mark_initiated",
			nil,
		).
			WithKind(richerror.KindConflict).
			WithCode(richerror.CodePaymentInvalidState)
	}

	payment.Status = status
	payment.ExternalID = externalID
	payment.PaymentURL = paymentURL

	f.markInitiatedCalls++

	return nil
}

type fakePaymentProvider struct {
	response paymentproviderparams.CreateResponse
	err      error
}

func (f *fakePaymentProvider) Create(_ context.Context, _ paymentproviderparams.CreateRequest) (paymentproviderparams.CreateResponse, error) {
	return f.response, f.err
}

func (f *fakePaymentProvider) Verify(context.Context, paymentproviderparams.VerifyRequest) (paymentproviderparams.VerifyResponse, error) {
	return paymentproviderparams.VerifyResponse{}, nil
}
