package telegramproduct

type starsPriceResponse struct {
	OK           bool   `json:"ok"`
	Kind         string `json:"kind"`
	PricePerStar string `json:"price_per_star"`
	MinAmount    uint32 `json:"min_amount"`
	MaxAmount    uint32 `json:"max_amount"`
}

type premiumResponse struct {
	OK    bool          `json:"ok"`
	Kind  string        `json:"kind"`
	Plans []premiumPlan `json:"plans"`
}

type premiumPlan struct {
	Months   int    `json:"months"`
	PriceUSD string `json:"price_usd"`
}
