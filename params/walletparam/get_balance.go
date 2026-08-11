package walletparam

import "telegram-service-platform/entity"

type GetBalanceRequest struct {
	UserID uint64
}

type GetBalanceResponse struct {
	Balance  entity.Amount
	Currency entity.Currency
}
