package smmentity

type Status struct {
	Charge     string `json:"charge"`
	StartCount string `json:"start_count"`
	Status     string `json:"status"`
	Remains    string `json:"remains"`
	Currency   string `json:"currency"`
}
