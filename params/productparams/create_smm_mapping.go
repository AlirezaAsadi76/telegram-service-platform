package productparams

import "telegram-service-platform/entity/smmentity"

type CreateSMMMappingRequest struct {
	SmmServiceId int64
	Name         string
	Platform     smmentity.PlatformType
	Category     smmentity.Category
	Description  string
	IsActive     bool
	ButtonName   string
}

type CreateSMMMappingResponse struct {
	ID           int64
	SmmServiceId int64
	Platform     smmentity.PlatformType
	Category     smmentity.Category
	ButtonName   string
}
