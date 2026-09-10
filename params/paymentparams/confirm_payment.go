package paymentparams

import "telegram-service-platform/entity/paymententity"

type ConfirmPaymentRequest struct {
	PaymentID    uint64
	ExternalID   string
	CallbackData map[string]any
}

type ConfirmPaymentResponse struct {
	PaymentID uint64
	OrderID   uint64
	Status    paymententity.PaymentStatus
}
