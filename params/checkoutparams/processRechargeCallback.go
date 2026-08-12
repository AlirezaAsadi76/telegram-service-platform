package checkoutparams

type ProcessRechargeCallbackRequest struct {
	UserId       uint64         `json:"user_id"`
	PaymentID    uint64         `json:"payment_id"`
	ExternalID   string         `json:"external_id"`
	CallbackData map[string]any `json:"callback_data"`
}
