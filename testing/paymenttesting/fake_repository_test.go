package paymenttesting

import (
	"context"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/pkg/richerror"
)

type fakePaymentRepository struct {
	payments             map[uint64]*paymententity.Payment
	paymentByIdempotency map[string]*paymententity.Payment

	createdPayment *paymententity.Payment

	createErr           error
	getByIDErr          error
	getByOrderIDErr     error
	getByIdempotencyErr error
	updateStatusErr     error
	markInitiatedErr    error

	getByIdempotencySequence []error

	getByIDCalls          int
	getByOrderIDCalls     int
	getByIdempotencyCalls int
	createCalls           int
	updateStatusCalls     int
	markInitiatedCalls    int
}

func newFakePaymentRepository() *fakePaymentRepository {
	return &fakePaymentRepository{
		payments:             make(map[uint64]*paymententity.Payment),
		paymentByIdempotency: make(map[string]*paymententity.Payment),
	}
}

func (f *fakePaymentRepository) Create(
	_ context.Context,
	payment *paymententity.Payment,
) error {
	f.createCalls++

	if f.createErr != nil {
		return f.createErr
	}

	payment.ID = 100

	f.createdPayment = payment
	f.payments[payment.ID] = payment
	f.paymentByIdempotency[payment.IdempotencyKey] = payment

	return nil
}

func (f *fakePaymentRepository) GetByID(
	_ context.Context,
	id uint64,
) (*paymententity.Payment, error) {
	f.getByIDCalls++

	if f.getByIDErr != nil {
		return nil, f.getByIDErr
	}

	payment, ok := f.payments[id]
	if !ok {
		return nil, richerror.New(
			"fakepayment.get_by_id",
			nil,
		).
			WithKind(richerror.KindNotFound).
			WithCode(richerror.CodePaymentNotFound)
	}

	return payment, nil
}
func (f *fakePaymentRepository) GetByExternalID(_ context.Context, external_id string) (*paymententity.Payment, error) {
	for _, payment := range f.payments {
		if payment.ExternalID == external_id {
			return payment, nil
		}
	}
	return nil, richerror.New("fakepayment.get_by_external_id", nil).
		WithKind(richerror.KindNotFound).
		WithCode(richerror.CodePaymentNotFound)
}
func (f *fakePaymentRepository) GetByOrderID(
	_ context.Context,
	orderID uint64,
) (*paymententity.Payment, error) {
	f.getByOrderIDCalls++

	if f.getByOrderIDErr != nil {
		return nil, f.getByOrderIDErr
	}

	for _, payment := range f.payments {
		if payment.OrderID == orderID {
			return payment, nil
		}
	}

	return nil, richerror.New(
		"fakepayment.get_by_order_id",
		nil,
	).
		WithKind(richerror.KindNotFound).
		WithCode(richerror.CodePaymentNotFound)
}

func (f *fakePaymentRepository) GetByIdempotencyKey(
	_ context.Context,
	key string,
) (*paymententity.Payment, error) {
	index := f.getByIdempotencyCalls
	f.getByIdempotencyCalls++

	if index < len(f.getByIdempotencySequence) {
		err := f.getByIdempotencySequence[index]
		if err != nil {
			return nil, err
		}
	}

	if f.getByIdempotencyErr != nil {
		return nil, f.getByIdempotencyErr
	}

	payment, ok := f.paymentByIdempotency[key]
	if !ok {
		return nil, richerror.New(
			"fakepayment.get_by_idempotency_key",
			nil,
		).
			WithKind(richerror.KindNotFound).
			WithCode(richerror.CodePaymentNotFound)
	}

	return payment, nil
}

func (f *fakePaymentRepository) UpdateStatus(
	_ context.Context,
	id uint64,
	status paymententity.PaymentStatus,
) error {
	f.updateStatusCalls++

	if f.updateStatusErr != nil {
		return f.updateStatusErr
	}

	payment, ok := f.payments[id]
	if !ok {
		return richerror.New(
			"fakepayment.update_status",
			nil,
		).
			WithKind(richerror.KindNotFound).
			WithCode(richerror.CodePaymentNotFound)
	}

	payment.Status = status
	return nil
}

func (f *fakePaymentRepository) MarkInitiated(
	_ context.Context,
	paymentID uint64,
	status paymententity.PaymentStatus,
	externalID string,
	paymentURL string,
) error {
	f.markInitiatedCalls++

	if f.markInitiatedErr != nil {
		return f.markInitiatedErr
	}

	payment, ok := f.payments[paymentID]
	if !ok {
		return richerror.New(
			"fakepayment.mark_initiated",
			nil,
		).
			WithKind(richerror.KindNotFound).
			WithCode(richerror.CodePaymentNotFound)
	}

	if payment.Status != paymententity.PaymentStatusCreating {
		return richerror.New(
			"fakepayment.mark_initiated",
			nil,
		).
			WithKind(richerror.KindConflict).
			WithCode(richerror.CodePaymentInvalidState)
	}

	payment.Status = status
	payment.ExternalID = externalID
	payment.PaymentURL = paymentURL

	return nil
}

func (f *fakePaymentRepository) GetPending(
	_ context.Context,
) ([]paymententity.Payment, error) {
	return nil, nil
}

func (f *fakePaymentRepository) GetExpired(
	_ context.Context,
) ([]paymententity.Payment, error) {
	return nil, nil
}
