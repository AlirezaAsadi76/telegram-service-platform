package orderentity

import "time"

type FulfillmentAttemptOutcome string

const (
	FulfillmentAttemptOutcomeCreated FulfillmentAttemptOutcome = "CREATED"
	FulfillmentAttemptOutcomeUnknown FulfillmentAttemptOutcome = "UNKNOWN"
)

type FulfillmentAttempt struct {
	ID              uint64
	OrderID         uint64
	ProviderID      uint64
	Outcome         FulfillmentAttemptOutcome
	ExternalOrderID string
	ResolvedAt      *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
