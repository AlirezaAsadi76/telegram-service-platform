package smmparams

type CreateOrderOutcome string

const (
	CreateOrderOutcomeCreated  CreateOrderOutcome = "CREATED"
	CreateOrderOutcomeRejected CreateOrderOutcome = "REJECTED"
	CreateOrderOutcomeUnknown  CreateOrderOutcome = "UNKNOWN"
)

type CreateOrderAdapterRequest struct {
	ProviderName string
	// ServiceID is the provider-specific service ID.
	ServiceID string
	Link      string
	Quantity  int64
}

type CreateOrderAdapterResponse struct {
	Outcome         CreateOrderOutcome
	ExternalOrderID string
}
