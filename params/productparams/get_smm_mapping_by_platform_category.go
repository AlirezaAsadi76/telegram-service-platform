package productparams

import "telegram-service-platform/entity/smmentity"

type GetSmmMappingByPlatformCategoryRequest struct {
	Platform smmentity.PlatformType
	Category smmentity.Category
}
type GetSmmMappingByPlatformCategoryResponse struct {
	SmmMapping []smmentity.SmmMapping
}
