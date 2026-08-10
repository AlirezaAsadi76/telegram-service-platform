package paymentparams

import (
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/paymententity"
)

type CreatePaymentRequest struct {
	OrderID  uint64
	UserID   uint64
	Method   paymententity.PaymentMethod
	Amount   float64
	Currency entity.Currency
}
type CreatePaymentResponse struct {
	PaymentID uint64
	Status    string
	Action    string
}
