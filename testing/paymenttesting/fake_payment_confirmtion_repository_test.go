package paymenttesting

import "context"

type fakePaymentConfirmationRepository struct {
	confirmErr              error
	failErr                 error
	markUnknownErr          error
	confirmCalls            int
	failCalls               int
	markUnknownCalls        int
	lastProviderReferenceID string
}

func newFakePaymentConfirmationRepository() *fakePaymentConfirmationRepository {
	return &fakePaymentConfirmationRepository{}
}
func (f *fakePaymentConfirmationRepository) Confirm(
	_ context.Context,
	_ uint64,
	providerReferenceID string,
) error {
	f.confirmCalls++
	f.lastProviderReferenceID = providerReferenceID
	return f.confirmErr
}

func (f *fakePaymentConfirmationRepository) Fail(
	_ context.Context,
	_ uint64,
) error {
	f.failCalls++
	return f.failErr
}

func (f *fakePaymentConfirmationRepository) MarkUnknown(
	_ context.Context,
	_ uint64,
) error {
	f.markUnknownCalls++
	return f.markUnknownErr
}
