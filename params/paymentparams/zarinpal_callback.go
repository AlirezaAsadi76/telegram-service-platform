package paymentparams

type ZarinpalCallback struct {
	Status    string `query:"status"`
	Authority string `query:"authority"`
}
