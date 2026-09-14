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

type StartPaymentHandlerRequest struct {
	OrderID        uint64                      `json:"order_id"`
	Method         paymententity.PaymentMethod `json:"method"`
	Amount         entity.Amount               `json:"amount"`
	Currency       entity.Currency             `json:"currency"`
	IdempotencyKey string                      `json:"idempotency_key"`
	CallbackURL    string                      `json:"callback_url"`
	Description    string                      `json:"description"`
}
