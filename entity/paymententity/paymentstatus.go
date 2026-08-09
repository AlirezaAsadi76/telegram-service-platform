package paymententity

type PaymentStatus string

const (
	PaymentPending    PaymentStatus = "PENDING"
	PaymentProcessing PaymentStatus = "PROCESSING"
	PaymentSuccess    PaymentStatus = "SUCCESS"
	PaymentFailed     PaymentStatus = "FAILED"
	PaymentCanceled   PaymentStatus = "CANCELED"
	PaymentExpired    PaymentStatus = "EXPIRED"
)
