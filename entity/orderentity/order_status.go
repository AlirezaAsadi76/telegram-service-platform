package orderentity

type OrderStatus string

const (
	OrderPending    OrderStatus = "PENDING"
	OrderProcessing OrderStatus = "PROCESSING"
	OrderSuccess    OrderStatus = "SUCCESS"
	OrderFailed     OrderStatus = "FAILED"
	OrderCanceled   OrderStatus = "CANCELED"
	OrderExpired    OrderStatus = "EXPIRED"
)
