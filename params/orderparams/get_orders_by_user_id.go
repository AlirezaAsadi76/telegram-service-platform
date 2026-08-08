package orderparams

type GetOrdersByUserIdRequest struct {
	UserId uint64
}
type GetOrdersByUserIdResponse struct {
	OrderInfo []OrderInfo
}
