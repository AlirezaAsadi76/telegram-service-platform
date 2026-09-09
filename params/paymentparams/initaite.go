package paymentparams

import "telegram-service-platform/entity/paymententity"

type InitiateRequest struct {
	PaymentID   uint64
	CallbackURL string
	Description string
}

type InitiateResponse struct {
	PaymentID  uint64
	Status     paymententity.PaymentStatus
	PaymentURL string
	ExternalID string
}
