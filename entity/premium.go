package entity

type PremiumPlan struct {
	Month    int64   `json:"month"`
	UsdPrice float64 `json:"usd_price"`
}
