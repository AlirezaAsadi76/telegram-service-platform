package paymentparams

type VerifyPaymentRequest struct {
	PaymentID   uint64
	ReferenceID string
}

type VerifyPaymentResponse struct {
}
