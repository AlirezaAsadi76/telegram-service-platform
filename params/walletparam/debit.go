package walletparam

import "telegram-service-platform/entity"

type DebitRequest struct {
	UserID         uint64
	Amount         entity.Amount
	ReferenceID    string
	IdempotencyKey string
}

type DebitResponse struct {
	TransactionID uint64
	NewBalance    entity.Amount
}
