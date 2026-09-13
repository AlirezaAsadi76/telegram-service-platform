package paymentflowtesting

import "telegram-service-platform/params/paymentparams"

type flowTestValidator struct {
	calls int
	err   error
}

func (f *flowTestValidator) ValidateZarinpalCallback(_ paymentparams.ZarinpalCallback) (map[string]string, error) {
	f.calls++

	return nil, f.err
}
