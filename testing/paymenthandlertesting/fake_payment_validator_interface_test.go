package paymenthandlertesting

import "telegram-service-platform/params/paymentparams"

type fakePaymentValidator struct {
	err                        error
	fieldErrs                  map[string]string
	calls                      int
	validateStartPaymentCalled bool
}

func (f *fakePaymentValidator) ValidateZarinpalCallback(_ paymentparams.ZarinpalCallback) (map[string]string, error) {
	f.calls++

	return f.fieldErrs, f.err
}
func (f *fakePaymentValidator) ValidateStartPayment(_ paymentparams.StartPaymentHandlerRequest) (map[string]string, error) {
	f.validateStartPaymentCalled = true

	return nil, nil
}
