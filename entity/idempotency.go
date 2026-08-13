package entity

type IdempotencyStatus string

const (
	IdempotencyStatusComplete   IdempotencyStatus = "complete"
	IdempotencyStatusPending    IdempotencyStatus = "pending"
	IdempotencyStatusError      IdempotencyStatus = "error"
	IdempotencyStatusProcessing IdempotencyStatus = "processing"
)

func (i IdempotencyStatus) String() string {
	return string(i)
}
