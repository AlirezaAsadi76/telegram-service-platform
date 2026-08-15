package smmentity

type SMM struct {
	Service  int64  `json:"service"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Rate     string `json:"rate"`
	Min      int64  `json:"min"`
	Max      int64  `json:"max"`
	Dripfeed bool   `json:"dripfeed"`
	Refill   bool   `json:"refill"`
	Cancel   bool   `json:"cancel"`
	Category string `json:"category"`
}
