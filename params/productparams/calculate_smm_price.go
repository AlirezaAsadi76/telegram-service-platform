package productparams

import (
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/productentity"
)

type CalculateSMMPriceRequest struct {
	MappingID int64
	Quantity  int64
}

type CalculateSMMPriceResponse struct {
	MappingID        int64
	ServiceID        int64
	Rate             entity.Amount
	PricePerThousand productentity.Price
	Price            productentity.Price
}
