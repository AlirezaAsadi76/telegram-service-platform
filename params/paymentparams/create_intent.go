package paymentparams

import (
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/paymententity"
)

type CreateIntentRequest struct {
	OrderID        uint64
	UserID         uint64
	Method         paymententity.PaymentMethod
	Amount         entity.Amount
	Currency       entity.Currency
	IdempotencyKey string
}

type CreateIntentResponse struct {
	PaymentID uint64
	Status    paymententity.PaymentStatus
}
