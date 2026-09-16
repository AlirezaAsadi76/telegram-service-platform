package smmprams

type CreateOrderOutcome string

const (
	CreateOrderOutcomeCreated  CreateOrderOutcome = "CREATED"
	CreateOrderOutcomeRejected CreateOrderOutcome = "REJECTED"
	CreateOrderOutcomeUnknown  CreateOrderOutcome = "UNKNOWN"
)

type CreateOrderAdapterRequest struct {
	ServiceID string
	Link      string
	Quantity  int64
}

type CreateOrderAdapterResponse struct {
	Outcome         CreateOrderOutcome
	ExternalOrderID string
}
