package zarinpal

type verifyResponse struct {
	Data   verifyData    `json:"data"`
	Errors *verifyErrors `json:"errors"`
}

type verifyData struct {
	Code        int64  `json:"code"`
	Message     string `json:"message"`
	CardHash    string `json:"card_hash"`
	CardPan     string `json:"card_pan"`
	RefID       int64  `json:"ref_id"`
	AgreementID string `json:"agreement_id"`
	FeeType     string `json:"fee_type"`
	Fee         int64  `json:"fee"`
}

type verifyErrors struct {
	Code        int64  `json:"code"`
	Message     string `json:"message"`
	Validations []any  `json:"validations"`
}
