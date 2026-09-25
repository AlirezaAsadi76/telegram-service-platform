package paymentparams

import (
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/paymententity"
	"time"
)

type CreateRequest struct {
	OrderID        uint64
	UserID         uint64
	Method         paymententity.PaymentMethod
	Amount         entity.Amount
	Currency       entity.Currency
	IdempotencyKey string
	ExpiredAt      time.Time
}

type CreateResponse struct {
	PaymentID  uint64
	PaymentURL string
	ExternalID string
}
