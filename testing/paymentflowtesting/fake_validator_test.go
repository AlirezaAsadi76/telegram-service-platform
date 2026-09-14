package paymentflowtesting

import "telegram-service-platform/params/paymentparams"

type flowTestValidator struct {
	calls int
	err   error
}

func (f *flowTestValidator) ValidateStartPayment(req paymentparams.StartPaymentHandlerRequest) (map[string]string, error) {

	return nil, nil
}

func (f *flowTestValidator) ValidateZarinpalCallback(_ paymentparams.ZarinpalCallback) (map[string]string, error) {
	f.calls++

	return nil, f.err
}
