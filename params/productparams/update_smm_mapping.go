package productparams

import "telegram-service-platform/entity/smmentity"

type UpdateSMMMappingRequest struct {
	Id           int64
	SmmServiceId int64
	Name         string
	Platform     smmentity.PlatformType
	Category     smmentity.Category
	Description  string
	IsActive     bool
	ButtonName   string
	SortOrder    int64
}

type UpdateSMMMappingResponse struct {
	ID           int64
	SmmServiceId int64
	Platform     smmentity.PlatformType
	Category     smmentity.Category
	ButtonName   string
	SortOrder    int64
}
