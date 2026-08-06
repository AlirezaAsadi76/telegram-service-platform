package exchangerate

type tonUsdPriceResponse map[string]price

type price struct {
	USD float64 `json:"usd"`
}

type marketsResponse struct {
	Result map[string]market `json:"result"`
}

type market struct {
	Symbol string `json:"symbol"`
	Stats  stats  `json:"stats"`
}

type stats struct {
	LastPrice float64 `json:"lastPrice"`
}
