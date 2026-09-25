package orderentity

import (
	"telegram-service-platform/entity"
)

type OrderFlowStage string

const (
	OrderFlowStageWaitingForLink     OrderFlowStage = "waiting_for_link"
	OrderFlowStageWaitingForQuantity OrderFlowStage = "waiting_for_quantity"
	OrderFlowStageConfirming         OrderFlowStage = "confirming"
	OrderFlowStageCompleted          OrderFlowStage = "completed"
)

type OrderFlowState struct {
	Stage                  OrderFlowStage  `json:"stage"`
	PurchaseID             string          `json:"purchase_id"`
	Platform               string          `json:"platform"`
	Category               string          `json:"category"`
	ServiceID              uint64          `json:"service_id"`
	ServiceName            string          `json:"service_name"`
	Rate                   entity.Amount   `json:"rate"`
	MinQuantity            int64           `json:"min_quantity"`
	MaxQuantity            int64           `json:"max_quantity"`
	Link                   string          `json:"link"`
	Quantity               int64           `json:"quantity"`
	Price                  entity.Amount   `json:"price"`
	Currency               entity.Currency `json:"currency"`
	PriceLockedAt          int64           `json:"price_locked_at"`
	PriceLockExpiresAt     int64           `json:"price_lock_expires_at"`
	PriceLockPaymentMethod string          `json:"price_lock_payment_method"`
	ExpiresAt              int64           `json:"expires_at"`
}
