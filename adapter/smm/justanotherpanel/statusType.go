package justanotherpanel

type StatusType string

const (
	statusCompleted    StatusType = "completed"
	statusCancelled    StatusType = "cancelled"
	statusFailed       StatusType = "failed"
	statusPending      StatusType = "pending"
	statusInProcessing StatusType = "in processing"
	statusProcessing   StatusType = "processing"
)
