package zarinpal

type VerifyRequest struct {
	MerchantID string `json:"merchant_id"`
	Authority  string `json:"authority"`
	Amount     int64  `json:"amount"`
}
