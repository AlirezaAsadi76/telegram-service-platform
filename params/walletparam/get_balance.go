package walletparam

import "telegram-service-platform/entity"

type GetBalanceRequest struct {
	UserID uint64
}

type GetBalanceResponse struct {
	Balance  entity.Amount   `json:"balance"`
	Currency entity.Currency `json:"currency"`
}
