package paymenthandler

type Handler struct {
	paymentService PaymentFlow
	paymentVal     PaymentValidator
}

func New(
	paymentService PaymentFlow,
	paymentVal PaymentValidator,
) *Handler {
	return &Handler{
		paymentService: paymentService,
		paymentVal:     paymentVal,
	}
}
