package zarinpal

type VerifyResponse struct {
	Data   VerifyData    `json:"data"`
	Errors *VerifyErrors `json:"errors"`
}

type VerifyData struct {
	Code        int64  `json:"code"`
	Message     string `json:"message"`
	CardHash    string `json:"card_hash"`
	CardPan     string `json:"card_pan"`
	RefID       int64  `json:"ref_id"`
	AgreementID string `json:"agreement_id"`
	FeeType     string `json:"fee_type"`
	Fee         int64  `json:"fee"`
}

type VerifyErrors struct {
	Code        int64  `json:"code"`
	Message     string `json:"message"`
	Validations []any  `json:"validations"`
}
