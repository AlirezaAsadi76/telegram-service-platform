package paymentflowtesting

import (
	"context"
	"telegram-service-platform/params/paymentproviderparams"
)

type flowTestProvider struct {
	createCalls int

	verifyCalls int

	lastVerifyRequest paymentproviderparams.VerifyRequest

	createResponse paymentproviderparams.CreateResponse
	createErr      error

	verifyResponse paymentproviderparams.VerifyResponse
	verifyErr      error
}

func (f *flowTestProvider) Create(_ context.Context, req paymentproviderparams.CreateRequest) (paymentproviderparams.CreateResponse, error) {
	f.createCalls++
	return f.createResponse, f.createErr
}

func (f *flowTestProvider) Verify(_ context.Context, req paymentproviderparams.VerifyRequest) (paymentproviderparams.VerifyResponse, error) {
	f.verifyCalls++
	f.lastVerifyRequest = req

	return f.verifyResponse, f.verifyErr
}
