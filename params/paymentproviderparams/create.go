package paymentproviderparams

import "telegram-service-platform/entity"

type CreateRequest struct {
	Amount   entity.Amount
	Currency entity.Currency
}

type CreateResponse struct {
	PaymentID  uint64
	PaymentURL string
	ExternalID string
}
