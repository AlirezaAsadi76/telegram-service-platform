package orderentity

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "PENDING"
	OrderStatusPaid       OrderStatus = "PAID"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusCompleted  OrderStatus = "COMPLETED"
	OrderStatusFailed     OrderStatus = "FAILED"
	OrderStatusCanceled   OrderStatus = "CANCELED"
	OrderStatusExpired    OrderStatus = "EXPIRED"
)
