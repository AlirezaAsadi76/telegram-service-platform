package smmprams

type RefillResponse struct {
	RefillID int64
}

type MultiRefillResponse struct {
	Refills map[string]int64 // order_id -> refill_id
}

type RefillStatusResponse struct {
	Status string
}

type MultiRefillStatusResponse struct {
	Statuses map[string]string // refill_id -> status
}
