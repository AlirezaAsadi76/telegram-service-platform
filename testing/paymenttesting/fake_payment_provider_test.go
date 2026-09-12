package paymenttesting

import (
	"context"

	"telegram-service-platform/params/paymentproviderparams"
)

type fakePaymentProvider struct {
	createResponse paymentproviderparams.CreateResponse
	createErr      error
	verifyResponse paymentproviderparams.VerifyResponse
	verifyErr      error

	createCalls int
	verifyCalls int

	lastCreateRequest paymentproviderparams.CreateRequest
	lastVerifyRequest paymentproviderparams.VerifyRequest
}

func (f *fakePaymentProvider) Create(
	_ context.Context,
	req paymentproviderparams.CreateRequest,
) (paymentproviderparams.CreateResponse, error) {
	f.createCalls++
	f.lastCreateRequest = req

	return f.createResponse, f.createErr
}

func (f *fakePaymentProvider) Verify(
	_ context.Context,
	req paymentproviderparams.VerifyRequest,
) (paymentproviderparams.VerifyResponse, error) {
	f.verifyCalls++
	f.lastVerifyRequest = req

	return f.verifyResponse, f.verifyErr
}
