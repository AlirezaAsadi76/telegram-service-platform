package paymentproviderparams

import "telegram-service-platform/entity/paymententity"

type VerifyRequest struct {
	PaymentID    uint64
	ExternalID   string
	CallbackData map[string]any
}

type VerifyResponse struct {
	Status paymententity.PaymentStatus
}
