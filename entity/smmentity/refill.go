package smmentity

// Refill represents a single refill status response from JAP.
type Refill struct {
	RefillID int64  `json:"refill"`
	Status   string `json:"status"`
	Error    string `json:"error"`
}

// MultiRefillItem represents one item in a multi-refill response.
// The "refill" field can be either an int64 (success) or {"error": "..."} (failure).
type MultiRefillItem struct {
	Order    int64  `json:"order"`
	RefillID int64  // parsed from raw refill
	Error    string // empty if success
}

// RefillStatusItem represents one item in a multi-refill-status response.
// The "status" field can be either a string (success) or {"error": "..."} (failure).
type RefillStatusItem struct {
	RefillID int64  `json:"refill"`
	Status   string // parsed from raw status
	Error    string // empty if success
}

type CancelItem struct {
	Order     int64  `json:"order"`
	Cancelled int64  // 1 = success, 0 = failed
	Error     string // empty if success
}
