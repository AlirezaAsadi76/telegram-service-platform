package productentity

import (
	"telegram-service-platform/entity"
)

type Price struct {
	USD   entity.Amount
	USDT  entity.Amount
	TON   entity.Amount
	Toman entity.Amount
}
