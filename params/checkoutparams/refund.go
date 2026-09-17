package checkoutparams

import "telegram-service-platform/entity"

type RefundOrderRequest struct {
	OrderID uint64
	Reason  string
}

type RefundOrderResponse struct {
	OrderID         uint64
	UserID          uint64
	WalletTxID      uint64
	RefundAmount    entity.Amount
	AlreadyRefunded bool
}
