package productparams

type GetStarsPriceResponse struct {
	OK           bool   `json:"ok"`
	Kind         string `json:"kind"`
	PricePerStar string `json:"price_per_star"`
	MinAmount    uint32 `json:"min_amount"`
	MaxAmount    uint32 `json:"max_amount"`
}
