package checkoutparams

import (
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/paymententity"
)

type ProcessWalletRechargeRequest struct {
	UserID   uint64
	Amount   entity.Amount
	Currency entity.Currency
	Method   paymententity.PaymentMethod
}

type ProcessWalletRechargeResponse struct {
	PaymentID  uint64
	PaymentURL string
}
