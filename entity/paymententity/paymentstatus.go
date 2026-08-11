package paymententity

type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "PENDING"
	PaymentStatusProcessing PaymentStatus = "PROCESSING"
	PaymentStatusSuccess    PaymentStatus = "SUCCESS"
	PaymentStatusFailed     PaymentStatus = "FAILED"
	PaymentStatusCanceled   PaymentStatus = "CANCELED"
	PaymentStatusExpired    PaymentStatus = "EXPIRED"
)
