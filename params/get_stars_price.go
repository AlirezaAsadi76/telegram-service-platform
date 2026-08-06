package params

type GetStarsPriceResponse struct {
	OK           bool   `json:"ok"`
	PricePerStar string `json:"price_per_star"`
}
