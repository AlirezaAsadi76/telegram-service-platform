package paymententity

type PaymentMethod string

const (
	PaymentMethodZarinpal PaymentMethod = "ZARINPAL"
	PaymentMethodCrypto   PaymentMethod = "CRYPTO"
	PaymentMethodManual   PaymentMethod = "MANUAL"
)
