package walletparam

import "telegram-service-platform/entity"

type CreditRequest struct {
	UserID         uint64
	Amount         entity.Amount
	ReferenceID    string
	IdempotencyKey string
}

type CreditResponse struct {
	TransactionID uint64
	NewBalance    entity.Amount
}
