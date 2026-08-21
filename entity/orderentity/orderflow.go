package orderentity

import (
	"telegram-service-platform/entity"
	"time"
)

type OrderFlowStage string

const (
	OrderFlowStageWaitingForLink     OrderFlowStage = "waiting_for_link"
	OrderFlowStageWaitingForQuantity OrderFlowStage = "waiting_for_quantity"
	OrderFlowStageConfirming         OrderFlowStage = "confirming"
	OrderFlowStageCompleted          OrderFlowStage = "completed"
)

// OrderFlowState وضعیت موقت سفارش کاربر در Redis
type OrderFlowState struct {
	Stage     OrderFlowStage  `json:"stage"`
	Platform  string          `json:"platform"`
	Category  string          `json:"category"`
	ServiceID uint64          `json:"service_id"`
	Link      string          `json:"link"`
	Quantity  int64           `json:"quantity"`
	Price     entity.Amount   `json:"price"`
	Currency  entity.Currency `json:"currency"`
	ExpiresAt time.Duration   `json:"expires_at"`
}
