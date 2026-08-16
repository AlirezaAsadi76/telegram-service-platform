package smmprams

import "telegram-service-platform/entity/smmentity"

type RefillResponse struct {
	RefillID int64
}

type MultiRefillResponse struct {
	Items []smmentity.MultiRefillItem
}

type RefillStatusResponse struct {
	Status string
	Error  string // empty if success
}

type MultiRefillStatusResponse struct {
	Items []smmentity.RefillStatusItem
}
