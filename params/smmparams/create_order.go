package smmparams

type CreateOrderResult struct {
	Outcome         CreateOrderOutcome
	ProviderID      uint64
	ProviderName    string
	ExternalOrderID string
}
