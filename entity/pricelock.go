package entity

type PriceLockPaymentMethod string

const (
	PriceLockPaymentMethodWallet   PriceLockPaymentMethod = "WALLET"
	PriceLockPaymentMethodZarinpal PriceLockPaymentMethod = "ZARINPAL"
	PriceLockPaymentMethodCrypto   PriceLockPaymentMethod = "CRYPTO"
)
