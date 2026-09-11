package paymentproviderparams

import (
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/paymententity"
)

type VerifyRequest struct {
	PaymentID    uint64
	ExternalID   string
	Amount       entity.Amount
	Currency     entity.Currency
	CallbackData map[string]any
}

type VerifyResponse struct {
	Status      paymententity.PaymentStatus
	ReferenceID string
}
