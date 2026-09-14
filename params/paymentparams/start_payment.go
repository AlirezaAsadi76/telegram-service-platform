package paymentparams

import (
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/paymententity"
)

type StartPaymentRequest struct {
	OrderID        uint64
	UserID         uint64
	Method         paymententity.PaymentMethod
	Amount         entity.Amount
	Currency       entity.Currency
	IdempotencyKey string
	CallbackURL    string
	Description    string
}
