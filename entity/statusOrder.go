package entity

type Status string

const (
	StatusPending        Status = "PENDING"
	StatusSuccess        Status = "SUCCESS"
	StatusError          Status = "ERROR"
	StatusPendingPayment Status = "PENDING_PAYMENT"
	StatusPaid           Status = "PAID"
	StatusProcessing     Status = "PROCESSING"
	StatusCompleted      Status = "COMPLETED"
	StatusFailed         Status = "FAILED"
	StatusCanceled       Status = "CANCELED"
)
