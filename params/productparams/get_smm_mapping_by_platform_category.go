package productparams

import "telegram-service-platform/entity/smmentity"

type GetSmmMappingByPlatformCategoryRequest struct {
	Platform smmentity.PlatformType
	Category smmentity.CategoryType
}
type GetSmmMappingByPlatformCategoryResponse struct {
	SmmMapping []smmentity.SmmMapping
}
