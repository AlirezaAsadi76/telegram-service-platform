package smmprams

import "telegram-service-platform/entity/smmentity"

type GetStatusResponse struct {
	Status smmentity.Status
}

type GetMultiStatusResponse struct {
	Statuses map[string]smmentity.Status
}
