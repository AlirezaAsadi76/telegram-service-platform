package paymenthandlertesting

import "telegram-service-platform/params/paymentparams"

type fakePaymentValidator struct {
	err       error
	fieldErrs map[string]string
	calls     int
}

func (f *fakePaymentValidator) ValidateZarinpalCallback(_ paymentparams.ZarinpalCallback) (map[string]string, error) {
	f.calls++

	return f.fieldErrs, f.err
}
