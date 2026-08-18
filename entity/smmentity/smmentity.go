package smmentity

type SMM struct {
	Id           int64  `json:"id"`
	Service      int64  `json:"service"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Rate         int64  `json:"rate"`
	Min          int64  `json:"min"`
	Max          int64  `json:"max"`
	DripFeed     bool   `json:"dripfeed"`
	Refill       bool   `json:"refill"`
	Cancel       bool   `json:"cancel"`
	IsActive     bool   `json:"is_active"`
	Category     string `json:"category"`
	ProviderName string `json:"provider_name"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type SmmMapping struct {
	Id           int64
	SmmServiceId int64
	Name         string
	Platform     PlatformType
	Category     Category
	Description  string
	IsActive     bool
	ButtonName   string
	CreatedAt    string
	UpdatedAt    string
}
