package paymenthandlertesting

import (
	"context"
	"telegram-service-platform/params/paymentparams"
)

type fakePaymentFlow struct {
	calls    int
	lastReq  paymentparams.ConfirmPaymentByExternalIDRequest
	response *paymentparams.ConfirmPaymentResponse
	err      error

	startPaymentCalled bool
	startPaymentReq    paymentparams.StartPaymentRequest

	startPaymentResponse *paymentparams.InitiateResponse
	startPaymentError    error
}

func (f *fakePaymentFlow) ConfirmPaymentByExternalID(_ context.Context, req paymentparams.ConfirmPaymentByExternalIDRequest) (*paymentparams.ConfirmPaymentResponse, error) {
	f.calls++
	f.lastReq = req

	return f.response, f.err
}

func (f *fakePaymentFlow) StartPayment(_ context.Context, req paymentparams.StartPaymentRequest) (*paymentparams.InitiateResponse, error) {
	f.startPaymentCalled = true
	f.startPaymentReq = req

	return f.startPaymentResponse, f.startPaymentError
}
