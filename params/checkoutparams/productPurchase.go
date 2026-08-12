package checkoutparams

import (
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/entity/productentity"
)

type ProcessProductPurchaseRequest struct {
	UserID      uint64
	ProductType productentity.ProductType
	ProductID   uint64
	Quantity    int64
	TargetLink  string
	Amount      entity.Amount
	Currency    entity.Currency
	Method      paymententity.PaymentMethod // WALLET, ZARINPAL, CRYPTO
}
