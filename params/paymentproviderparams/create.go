package paymentproviderparams

import "telegram-service-platform/entity"

type CreateRequest struct {
	PaymentID   uint64
	OrderID     uint64
	Amount      entity.Amount
	Currency    entity.Currency
	CallbackURL string
	Description string
}

type CreateResponse struct {
	ExternalID string
	PaymentURL string
}
