package productparams

import "telegram-service-platform/entity/smmentity"

type GetSmmMappingByIDRequest struct {
	Id int64
}

type GetSmmMappingByIDResponse struct {
	SmmMapping *smmentity.SmmMapping
}
