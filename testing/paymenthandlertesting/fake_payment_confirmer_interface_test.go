package paymenthandlertesting

import (
	"context"
	"telegram-service-platform/params/paymentparams"
)

type fakePaymentConfirmer struct {
	calls    int
	lastReq  paymentparams.ConfirmPaymentByExternalIDRequest
	response *paymentparams.ConfirmPaymentResponse
	err      error
}

func (f *fakePaymentConfirmer) ConfirmPaymentByExternalID(_ context.Context, req paymentparams.ConfirmPaymentByExternalIDRequest) (*paymentparams.ConfirmPaymentResponse, error) {
	f.calls++
	f.lastReq = req

	return f.response, f.err
}
