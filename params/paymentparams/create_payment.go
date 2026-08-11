package paymentparams

import (
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/paymententity"
)

type CreateRequest struct {
	OrderID        *uint64
	UserID         uint64
	Method         paymententity.PaymentMethod
	Amount         entity.Amount
	Currency       entity.Currency
	IdempotencyKey string
	ExpiredAt      *interface{}
}

type CreateResponse struct {
	PaymentID  uint64
	PaymentURL string
	ExternalID string
}
