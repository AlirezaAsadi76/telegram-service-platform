package paymenthandler

import "telegram-service-platform/validator/paymentvalidator"

type Handler struct {
	paymentService PaymentConfirmer
	paymentVal     paymentvalidator.Validator
}

func New(
	paymentService PaymentConfirmer,
	paymentVal paymentvalidator.Validator,
) *Handler {
	return &Handler{
		paymentService: paymentService,
		paymentVal:     paymentVal,
	}
}
