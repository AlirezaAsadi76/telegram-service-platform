package paymententity

type PaymentMethod string

const (
	PaymentMethodGateway PaymentMethod = "GATEWAY"
	PaymentMethodCrypto  PaymentMethod = "CRYPTO"
	PaymentMethodManual  PaymentMethod = "MANUAL"
)
