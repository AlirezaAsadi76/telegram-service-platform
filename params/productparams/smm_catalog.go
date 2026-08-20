package productparams

import "telegram-service-platform/entity/smmentity"

type GetDistinctPlatformsRequest struct {
}

type GetDistinctPlatformsResponse struct {
	Platforms []smmentity.Platform
}

type GetDistinctCategoriesByPlatformRequest struct {
	Platform smmentity.PlatformType
}

type GetDistinctCategoriesByPlatformResponse struct {
	Categories []smmentity.Category
}
