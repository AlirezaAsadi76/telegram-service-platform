package orderentity

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "PENDING"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusSuccess    OrderStatus = "SUCCESS"
	OrderStatusFailed     OrderStatus = "FAILED"
	OrderStatusCanceled   OrderStatus = "CANCELED"
	OrderStatusExpired    OrderStatus = "EXPIRED"
)
