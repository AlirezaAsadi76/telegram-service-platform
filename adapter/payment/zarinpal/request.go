package zarinpal

type VerifyRequest struct {
	MerchantID string `json:"merchant_id"`
	Authority  string `json:"authority"`
	Amount     int64  `json:"amount"`
}

type createRequest struct {
	MerchantID  string `json:"merchant_id"`
	Amount      int64  `json:"amount"`
	CallbackURL string `json:"callback_url"`
	Description string `json:"description"`
}
