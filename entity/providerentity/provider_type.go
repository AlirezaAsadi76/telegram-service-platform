package providerentity

type ProviderType string

// SMM / PAYMENT / EXCHANGE
const (
	ProviderTypeSMM      ProviderType = "SMM"
	ProviderTypePayment  ProviderType = "PAYMENT"
	ProviderTypeExchange ProviderType = "EXCHANGE"
)
