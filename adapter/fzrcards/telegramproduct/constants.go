package telegramproduct

import (
	"net/http"
)

const (
	getStarsPricePath = "telegram/stars"
	getPremiumPath    = "telegram/premium"
)

const (
	getStarsPriceMethod = http.MethodGet
	getPremiumMethod    = http.MethodGet
)
